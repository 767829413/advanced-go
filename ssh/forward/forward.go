package forward

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

type TargetServers struct {
	LocalPort  string
	TargetAddr string
	TargetPort string
}

func SSHWorkForward(
	useName, target, port, keyPath string,
	targetServers []*TargetServers,
) (*ssh.Client, error) {
	// 读取私钥文件
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %v", err)
	}

	// 创建 signer
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
	}

	// 配置 SSH 客户端
	config := &ssh.ClientConfig{
		User: useName,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	// 连接到远程服务器
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", target, port), config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
	}

	// 设置本地端口转发
	for _, targetServer := range targetServers {
		go getPortForwardingFunc(
			client,
			targetServer.LocalPort,
			targetServer.TargetAddr,
			targetServer.TargetPort,
		)()
	}

	return client, nil
}

func getPortForwardingFunc(client *ssh.Client, localPort, targetAddr, targetPort string) func() {
	return func() {
		localListener, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", localPort))
		if err != nil {
			fmt.Printf("Failed to listen on local port %s: %v\n", localPort, err)
			return
		}
		defer localListener.Close()

		for {
			localConn, err := localListener.Accept()
			if err != nil {
				fmt.Printf("Failed to accept connection on local port: %v\n", err)
				continue
			}

			remoteConn, err := client.Dial(
				"tcp",
				fmt.Sprintf("%s:%s", targetAddr, targetPort),
			)
			if err != nil {
				fmt.Printf("Failed to connect to remote server: %v\n", err)
				localConn.Close()
				continue
			}

			go forwardConnection(localConn, remoteConn)
		}
	}
}

func forwardConnection(local, remote net.Conn) {
	defer local.Close()
	defer remote.Close()

	done := make(chan bool, 2)

	go func() {
		_, _ = io.Copy(local, remote)
		done <- true
	}()

	go func() {
		_, _ = io.Copy(remote, local)
		done <- true
	}()

	<-done
}
