package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	webhookURL       = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="
	SendTextType     = "text"
	SendMarkdownType = "markdown"
)

type message struct {
	MsgType  string       `json:"msgtype"`
	Text     *textMsg     `json:"text,omitempty"`
	Markdown *markdownMsg `json:"markdown,omitempty"`
	// 可以根据需要添加其他类型的消息结构
}

type textMsg struct {
	Content             string   `json:"content"`
	MentionedList       []string `json:"mentioned_list,omitempty"`
	MentionedMobileList []string `json:"mentioned_mobile_list,omitempty"`
}

type markdownMsg struct {
	Content string `json:"content"`
}

func SendMessage(key, msgType, content string, mentionAll bool) error {
	msg := message{
		MsgType: msgType,
	}

	switch msgType {
	case "text":
		msg.Text = &textMsg{
			Content: content,
		}
		if mentionAll {
			msg.Text.MentionedList = []string{"@all"}
		}
	case "markdown":
		msg.Markdown = &markdownMsg{
			Content: content,
		}
		if mentionAll {
			msg.Markdown.Content += "\n<@all>"
		}
	default:
		return fmt.Errorf("unsupported message type: %s", msgType)
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("JSON编码失败: %v", err)
	}

	resp, err := http.Post(webhookURL+key, "application/json", bytes.NewBuffer(jsonMsg))
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()
	return nil
}
