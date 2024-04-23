package notify

import (
	"encoding/json"
	"fmt"
	httpUtil "github.com/767829413/advanced-go/util/http"
)

// 通知消息
// 通知对象
func SendNotify(notifyStr string, rootKey string) {
	if len(rootKey) == 0 {
		return
	}
	defer func() {
		if r := recover(); r != nil {

		}
	}()
	url := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=" + rootKey
	body := httpUtil.WithBody(getMarkDownMsg(notifyStr))
	res, e := httpUtil.Post(url, body)
	fmt.Println(res, e)
}

type msgInfo struct {
	Msgtype  string      `json:"msgtype"`
	Markdown interface{} `json:"markdown"`
}

type markDownMsg struct {
	Content interface{} `json:"content"`
}

func getMarkDownMsg(msg string) []byte {
	t := &msgInfo{
		Msgtype:  "markdown",
		Markdown: markDownMsg{Content: msg},
	}
	res, _ := json.Marshal(t)
	return res
}
