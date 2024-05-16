package main

import (
	"fmt"
	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
	"log"
)

func main() {

	// Start new ssh connection with private key.
	auth, err := goph.Key("C:/Users/NUC/.ssh/id_rsa_wsl", "")
	if err != nil {
		log.Fatal("goph.Key error", err)
	}

	client, err := goph.NewConn(&goph.Config{
		User:     "fangyuan",
		Addr:     "116.62.116.90",
		Port:     10022,
		Auth:     auth,
		Timeout:  goph.DefaultTimeout,
		Callback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		log.Fatal("goph.NewConn error", err)
	}

	// Defer closing the network connection.
	defer client.Close()

	// Execute your command.
	out, err := client.Run(
		`sh totalLog.sh`,
	)

	if err != nil {
		log.Fatal("client.Run: ", err)
	}

	// Get your output as []byte.
	fmt.Println(string(out))
}
