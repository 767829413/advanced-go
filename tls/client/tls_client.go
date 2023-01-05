package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"sync"

	"go.uber.org/ratelimit"
)

func main() {
	url := "https://localhost:4443"
	connNum := 50000
	qps := 3000

	bucket := ratelimit.New(int(qps))

	var l sync.Mutex
	connList := make([]*http.Client, connNum)

	for i := 0; ; i++ {
		bucket.Take()
		i := i
		go func() {
			l.Lock()
			if connList[i%len(connList)] == nil {
				connList[i%len(connList)] = &http.Client{
					Transport: &http.Transport{
						TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
						IdleConnTimeout:     0,
						MaxIdleConns:        1,
						MaxIdleConnsPerHost: 1,
					},
				}
			}
			conn := connList[i%len(connList)]
			l.Unlock()
			if resp, e := conn.Get(url); e != nil {
				fmt.Println(e)
			} else {
				defer resp.Body.Close()
				io.ReadAll(resp.Body)
			}
		}()
	}
}
