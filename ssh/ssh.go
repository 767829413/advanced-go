package main

import (
	"fmt"
	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 4 {
		log.Fatal("Usage: go run main.go <user> <ip> <port>")
	}

	user := os.Args[1]
	ip := os.Args[2]
	port, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Fatalf("Invalid port number: %s", os.Args[3])
	}

	// Start new ssh connection with private key.
	auth, err := goph.Key("C:/Users/NUC/.ssh/id_rsa_wsl", "")
	if err != nil {
		log.Fatal("goph.Key error", err)
	}

	client, err := goph.NewConn(&goph.Config{
		User:     user,
		Addr:     ip,
		Port:     uint(port),
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
