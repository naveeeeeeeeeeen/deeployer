package tables

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

func (c Config) SSHConnect() {
	keyPath := "/Users/naveenwork/.ssh/id_ed25519"
	command := "ls /home/naveen"

	key, err := os.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("Cant read file %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)

	if err != nil {
		log.Fatalf("cant parse key, %v", err)
	}

	config := &ssh.ClientConfig{
		User: c.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	conn, err := ssh.Dial("tcp", c.Host+":22", config)
	if err != nil {
		log.Fatalf("err connecting via ssh %v", err)
	}

	defer conn.Close()

	session, err := conn.NewSession()

	if err != nil {
		log.Fatalf("unaable to create connection, %v", err)
	}

	defer session.Close()

	output, err := session.CombinedOutput(command)
	if err != nil {
		log.Fatalf("command exection failed, err: %v", err)
	}

	fmt.Printf("Output: \n%s", output)

}
