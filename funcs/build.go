package funcs

import (
	"deeployer/tables"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func cloneRepo(repoUrl string, accessToken string) string {
	tmpDir, err := os.MkdirTemp("", "build-repos")
	auth := &http.BasicAuth{
		Username: "git",
		Password: accessToken,
	}

	if err != nil {
		log.Fatalf("Failed to create repo, %v", err)
	}
	fmt.Println("Cloning repo ", repoUrl)
	_, err = git.PlainClone(tmpDir, false, &git.CloneOptions{
		URL:           repoUrl,
		Progress:      os.Stdout,
		Auth:          auth,
		ReferenceName: "master",
	})
	if err != nil {
		log.Fatalf("error cloning repo, err: %v", err)
	}
	fmt.Println("Clone the repo successfully, dir: ", tmpDir)
	return tmpDir
}

func uploadDir(client *sftp.Client, localPath string, remotePath string) error {
	entries, err := os.ReadDir(localPath)
	if err != nil {
		return err
	}

	err = client.MkdirAll(remotePath)

	if err != nil {
		fmt.Println("error creating directory")
		return err
	}

	for _, entry := range entries {
		localFile := filepath.Join(localPath, entry.Name())
		remoteFile := filepath.Join(remotePath, entry.Name())

		if entry.IsDir() {
			err = uploadDir(client, localFile, remoteFile)
			if err != nil {
				return err
			}
		} else {
			src, err := os.Open(localFile)
			if err != nil {
				return err
			}
			defer src.Close()

			dest, err := client.Create(remoteFile)
			if err != nil {
				return err
			}
			defer dest.Close()

			_, err = io.Copy(dest, src)

			if err != nil {
				return err
			}
		}
	}
	return nil
}

func sendBuildFiles(dir string, c *tables.Config) error {
	key, err := os.ReadFile("/Users/naveenwork/.ssh/ec2.pem")
	if err != nil {
		return err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return err
	}

	config := &ssh.ClientConfig{
		User: "ubuntu",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	conn, err := ssh.Dial("tcp", c.Host+":22", config)
	defer conn.Close()
	client, err := sftp.NewClient(conn)
	if err != nil {
		return err
	}
	log.Println("opening directory for inspection")
	dest, err := os.Open(dir + "/" + c.BuildFilePath.String)
	defer dest.Close()
	if err != nil {
		return err
	}
	pathInfo, err := dest.Stat()
	if err != nil {
		return err
	}
	if pathInfo.IsDir() {
		err = uploadDir(client, dir+"/"+c.BuildFilePath.String, c.DestPath.String)
	} else {
		// upload just a single build file
		log.Println(dir + "/" + c.BuildFilePath.String)
		log.Println(c.DestPath.String + "/" + c.BuildFilePath.String)
		src, err := os.Open(dir + "/" + c.BuildFilePath.String)
		if err != nil {
			return err
		}
		defer src.Close()
		_, err = client.Stat(c.DestPath.String)
		if err != nil {
			log.Println("creating new dir in remote")
			err = client.Mkdir(c.DestPath.String)
			if err != nil {
				log.Println("error creating a dir", err)
			}
		}
		d, err := client.Create(c.DestPath.String + "/" + c.BuildFilePath.String)
		if err != nil {
			return err
		}
		defer d.Close()

		_, err = io.Copy(d, src)
		if err != nil {
			return err
		}
		log.Println("uploaded file successfully")
	}

	var moveFiles bool = false
	if err != nil {
		log.Println("error while uploading Dir ", err)
		return err
	}

	if moveFiles {
		// SINCE WE REQUIRE ROOT PERMISSIONS TO COPY FILE TO /var/http we are better of copy it to home
		// and thne move them to /var/http
		session, err := conn.NewSession()
		if err != nil {
			return err
		}
		command := "sudo rm -rf /var/http/* && sudo mv /home/ubuntu/dist/* /var/http/"
		session.Stdin = os.Stdin
		session.Stdout = os.Stdout
		session.Stderr = os.Stderr
		return session.Run(command)
	}
	return nil
}

func runBuildCommands(commands string) string {
	out, err := exec.Command("bash", "-c", commands).Output()
	if err != nil {
		log.Fatalf("error running commands, err: %v", out)
	}
	return string(out)
}

func Build(configId string) error {
	c, err := tables.GetProjectConfig(configId)
	if err != nil {
		return fmt.Errorf("error getting config, %v ", err)
	}
	dir := cloneRepo(c.RepoUrl, c.GitKey)
	commands := c.BuildCommands.String

	defer os.RemoveAll(dir)
	buildCommand := "cd " + dir + " && " + commands
	buildOuput := runBuildCommands(buildCommand)
	log.Println("Build output: \n", buildOuput)
	err = sendBuildFiles(dir, &c)
	if err != nil {
		log.Fatalf("Error uploading file, %v", err)
	}
	log.Println("Build successful")
	return err
}
