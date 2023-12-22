package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
)

type QueryParams map[string]string

type HwSsoInfo struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    any    `json:"expires_in"`
	Scope        string `json:"scope"`
	Id           string `json:"id"`
	Type         any    `json:"type"`
	Mobile       string `json:"mobile"`
}

func GetHwRequsetParam(clientId, code string) (requestParam map[string]string, url string) {
	requestParam = map[string]string{
		"client_id":     clientId,
		"client_secret": "xxxxxxx",
		"code":          code,
		"grant_type":    "authorization_code",
	}
	url = "xxxx"
	return
}

func main() {
	requestParam, url := GetHwRequsetParam("474A8791134B39", "Qoq9at")
	res := HttpPost[HwSsoInfo](
		url,
		requestParam,
	)
	fmt.Println(res)
}

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func HttpPost[T any](url string, params map[string]string) *T {
	client := resty.New()
	reqClient := client.R()
	reqClient.SetHeader("Content-Type", "multipart/form-data;")
	reqClient.SetFormData(params)
	if params != nil {
		reqClient.SetBody(params)
	}
	resp, err := reqClient.Post(url)

	if err != nil {
		log.Panic("resp Error: %s\n", err.Error())
		return nil
	}

	log.Println("Response Body: %s\n", string(resp.Body()))
	// 对body进行解析
	var response T
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		log.Panic("json Error: %s\n", err.Error())
		return nil
	}
	return &response
}
