// +build !windows

package main

import (
	"log"
	"net"
	"os"
)

func dialAgent() (net.Conn, error) {
	sock := os.Getenv("SSH_AUTH_SOCK")
	if sock == "" {
		log.Fatal("SSH_AUTH_SOCK is not set")
	}
	return net.Dial("unix", sock)
}
