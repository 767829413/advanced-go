package util

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	logger "log"
	"reflect"
	"sort"
	"strings"
)

type QueryParams map[string]string

func StructToQueryParams(s interface{}) QueryParams {
	data := make(map[string]string)

	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()

		// 忽略未导出的字段
		if field.PkgPath != "" {
			continue
		}

		// 转换为字符串并添加到 map 中
		data[field.Name] = fmt.Sprintf("%v", value)
	}

	return data
}

func HttpPost[T any](url string, query QueryParams, params map[string]any) *T {
	client := resty.New()
	reqClient := client.R()
	if query != nil {
		reqClient.SetQueryParams(query)
	}
	if params != nil {
		reqClient.SetBody(params)
	}
	resp, err := reqClient.Post(url)

	if err != nil {
		logger.Fatalf("resp Error: %s\n", err.Error())
		return nil
	}

	logger.Printf("Response Body: %s\n", string(resp.Body()))
	// 对body进行解析
	var response T
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		logger.Fatalf("json Error: %s\n", err.Error())
		return nil
	}
	return &response
}

func HttpPostByFrom[T any](url string, params map[string]string) (*T, error) {
	client := resty.New()
	reqClient := client.R()
	reqClient.SetHeader("Content-Type", "multipart/form-data;")
	reqClient.SetFormData(params)

	if params != nil {
		reqClient.SetBody(params)
	}
	resp, err := reqClient.Post(url)

	if err != nil {
		logger.Fatalf("resp Error: %s\n", err.Error())
		return nil, err
	}

	logger.Printf("Response Body: %s\n", string(resp.Body()))
	// 判断是否请求失败
	switch resp.StatusCode() {
	case 200:
		// 对body进行解析
		var response T
		err = json.Unmarshal(resp.Body(), &response)
		if err != nil {
			logger.Fatalf("json Error: %s\n", err.Error())
			return nil, err
		}
		return &response, nil
	default:
		return nil, errors.New(string(resp.Body()))
	}
}

func OpenEncrypt(appSecret string, params map[string]string) (string, error) {
	if len(appSecret) == 0 {
		return "", fmt.Errorf("miss appSecret")
	}
	keys := make([]string, 0, len(params))
	for k, _ := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	keyValues := []string{}
	for _, k := range keys {
		if k == "signature" || len(k) == 0 {
			continue
		}
		keyValues = append(keyValues, k+"="+params[k])
	}
	p := strings.Join(keyValues, "&")
	mac := hmac.New(sha1.New, []byte(appSecret))
	mac.Write([]byte(p))
	var signature = fmt.Sprintf("%X", mac.Sum(nil))
	return signature, nil
}
