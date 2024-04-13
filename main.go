package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	// etcd 集群化+客户端多Endpoints
	url := "http://192.168.20.20/manage/nc/login/doLoginForManager" // 替换为您要请求的POST接口URL

	// 设置定时器，每隔一秒发送一次POST请求
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// 构造POST请求的body
		jsonStr := []byte(`{
			"loginName": "t2zwt2",
			"password": "f379eaf3c831b04de153469d1bec345e",
			"loginType": 0
			
		}`) // 替换为您要发送的JSON数据

		// 发送POST请求
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		log.Println("Response Status:", resp.Status)
		// 读取并打印响应内容
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading response body:", err)
			continue
		}
		log.Println("Response Body:", string(body))
		resp.Body.Close()
	}
}
