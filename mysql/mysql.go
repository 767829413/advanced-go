package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {

	gorm.Open(mysql.Open(""), &gorm.Config{
		Logger: logger.Default,
	})

	
}
