// +build windows

package main

import (
	"net"
	"os"

	"github.com/Microsoft/go-winio"
)

func dialAgent() (net.Conn, error) {
	sock := os.Getenv("SSH_AUTH_SOCK")
	if sock == "" {
		sock = `\\.\pipe\openssh-ssh-agent`
	}
	return winio.DialPipe(sock, nil)
}
