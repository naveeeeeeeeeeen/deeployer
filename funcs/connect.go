package funcs

import "fmt"

// connect via ssh to the remote server

type Connection interface {
	SSHConnect()
}

func Connect(c Connection) {
	fmt.Println("Connecting via ssh")
	c.SSHConnect()
}
