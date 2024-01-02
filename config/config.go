package config

import (
	"encoding/json"
	"log"
	"os"
)

var (
	config conf
	path   = "D:/Document/code/go/src/github.com/767829413/advanced-go/config/conf/conf.json"
)

func init() {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Panic("无法读取配置文件:", err)
		return
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Panic("无法解析配置文件:", err)
		return
	}
}

type conf struct {
	MysqlUrl string `json:"mysqlUrl"` // mysql数据库连接url
}

func GetConfig() conf {
	return config
}
