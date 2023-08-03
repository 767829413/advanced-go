package main

import (
	"fmt"
	"net"
	"time"

	"github.com/767829413/advanced-go/util"
)

func main() {
	conn, err := net.Dial("udp", "localhost:9899")
	if err != nil {
		fmt.Println("net.Dial error: ", err)
		return
	}
	defer conn.Close()
	ticker := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-ticker.C:
			time.Sleep(1 * time.Second)
		default:
			fmt.Println("start request server")
			_, err = conn.Write([]byte(util.RandStr(5)))
			if err != nil {
				fmt.Println("conn.Write failed: ", err)
				return
			}
			buf := make([]byte, 4096)
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println("conn.Read failed: ", err)
			}
			fmt.Println("conn.Read: ", string(buf[:n]))
		}
	}
}
