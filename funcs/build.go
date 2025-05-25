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

func sendBuildFiles(dir string, host string) error {
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

	conn, err := ssh.Dial("tcp", host+":22", config)
	defer conn.Close()
	client, err := sftp.NewClient(conn)
	if err != nil {
		return err
	}
	// SINCE WE REQUIRE ROOT PERMISSIONS TO COPY FILE TO /var/http we are better of copy it to home
	// and thne move them to /var/http
	err = uploadDir(client, dir+"/exc", "/home/ubuntu/deeployer/")

	if err != nil {
		return err
	}

	return nil
	/*
		session, err := conn.NewSession()
		if err != nil {
			return err
		}
			command := "sudo rm -rf /var/http/* && sudo mv /home/ubuntu/dist/* /var/http/"
			session.Stdin = os.Stdin
			session.Stdout = os.Stdout
			session.Stderr = os.Stderr
			return session.Run(command)
	*/
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
	var buildCommands string

	if c.BuildCommands.Valid {
		buildCommands = c.BuildCommands.String
	} else {
		buildCommands = "bash && cd " + dir + " && npm install && npm run build"
	}
	defer os.RemoveAll(dir)

	buildOuput := runBuildCommands(buildCommands)
	log.Println("Build output: \n", buildOuput)
	err = sendBuildFiles(dir, c.Host)
	if err != nil {
		log.Fatalf("Error uploading file, %v", err)
	}
	return err
}
