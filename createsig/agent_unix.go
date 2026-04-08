// +build !windows

package main

import (
	"fmt"
	"net"
	"os"
)

func dialAgent() (net.Conn, error) {
	sock := os.Getenv("SSH_AUTH_SOCK")
	if sock == "" {
		return nil, fmt.Errorf("SSH_AUTH_SOCK is not set")
	}
	return net.Dial("unix", sock)
}
