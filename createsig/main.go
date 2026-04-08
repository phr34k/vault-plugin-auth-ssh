package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func genSig(nonce, privatekey, password string) {
	var (
		signer ssh.Signer
		err    error
	)

	if privatekey == "agent" || strings.HasPrefix(privatekey, "agent:") || privatekey == "" {
		sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
		if err != nil {
			log.Fatal(err)
		}
		defer sshAgent.Close()
		ag := agent.NewClient(sshAgent)
		signers, err := ag.Signers()
		if err != nil {
			log.Fatal(err)
		}
		if len(signers) == 0 {
			log.Fatal("no keys in agent")
		}

		if privatekey == "agent" || privatekey == "" {
			signer = signers[0]
		} else {
			keys, err := ag.List()
			if err != nil {
				log.Fatal(err)
			}
			target := strings.TrimPrefix(privatekey, "agent:")
			for i, s := range signers {
				if keys[i].Comment == target || ssh.FingerprintSHA256(s.PublicKey()) == target {
					signer = s
					break
				}
			}
			if signer == nil {
				log.Fatalf("key %s not found in agent", target)
			}
		}
	} else {
		pemBytes, err := os.ReadFile(privatekey)
		if err != nil {
			log.Fatal(err)
		}

		if password == "" {
			signer, err = ssh.ParsePrivateKey(pemBytes)
			if err != nil {
				log.Fatalf("parse key failed:%v", err)
			}
		} else {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(password))
			if err != nil {
				log.Fatalf("parse key failed:%v", err)
			}
		}
	}

	var signBytes []byte

	if nonce == "" {
		t := time.Now()
		timeBytes, _ := t.MarshalBinary()
		signBytes = append(signBytes, timeBytes...)
	} else {
		signBytes = append(signBytes, []byte(nonce)...)
	}

	res, err := signer.Sign(rand.Reader, signBytes)
	if err != nil {
		log.Fatal(err)
	}

	signatureBlob := res.Blob

	fmt.Println("signature=" + base64.StdEncoding.EncodeToString(signatureBlob) + " nonce=" + base64.StdEncoding.EncodeToString(signBytes))
}

func printHelp() {
	fmt.Println("This tool will print out a signature based on a nonce to be used with vault-plugin-auth-ssh")
	fmt.Println("You can get a nonce by running \"vault read auth/ssh/nonce\"")
	fmt.Println("")
	fmt.Println("Need " + os.Args[0] + " <nonce> <key-path> <password>")
	fmt.Println("eg. " + os.Args[0] + " anonce ~/.ssh/id_rsa mypassword")
	fmt.Println("")
	fmt.Println("To use the ssh-agent, use \"agent\" as <key-path> or omit it")
	fmt.Println("To use a specific key from the agent, use \"agent:<comment>\" or \"agent:<fingerprint>\"")
	fmt.Println("eg. " + os.Args[0] + " anonce agent:mykey")
	fmt.Println("eg. " + os.Args[0] + " anonce agent:SHA256:pG9... ")
	fmt.Println("eg. " + os.Args[0] + " anonce agent")
	fmt.Println("eg. " + os.Args[0] + " anonce")
	fmt.Println("")
	fmt.Println("If you don't have a password just omit it")
	fmt.Println("eg. " + os.Args[0] + " anonce ~/.ssh/id_rsa")
}

func main() {
	switch len(os.Args) {
	case 2:
		genSig(os.Args[1], "", "")
	case 3:
		genSig(os.Args[1], os.Args[2], "")
	case 4:
		genSig(os.Args[1], os.Args[2], os.Args[3])
	default:
		printHelp()
	}
}
