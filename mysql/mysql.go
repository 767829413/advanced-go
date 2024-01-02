package main

import (
	"github.com/767829413/advanced-go/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

func main() {
	conf := config.GetConfig()
	db, err := gorm.Open(mysql.Open(conf.MysqlUrl), &gorm.Config{
		Logger: logger.Default,
	})
	if err != nil {
		log.Panic("gorm.Open error: ", err)
	}
	var res int
	db.Raw("select 1").Scan(&res)
	log.Println("res: ", res)
}
