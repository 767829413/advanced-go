package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

type bodyParam struct {
	Url string `json:"url"`
}

// MarshalJSON 实现 json.Marshaler 接口
func (bp bodyParam) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"url": json.RawMessage(fmt.Sprintf(`"%s"`, bp.Url)),
	})
}

func main() {
	client := resty.New()

	body := bodyParam{
		Url: "",
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post("https://your-api-endpoint.com")

	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	fmt.Println("Response Status Code:", resp.StatusCode())
	fmt.Println("Response Body:", string(resp.Body()))
}
