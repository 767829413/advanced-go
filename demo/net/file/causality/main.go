package main

import (
	"flag"
	"fmt"
	"os"
)

var path string

func init() {
	flag.StringVar(&path, "path", "", "file path")
	flag.Parse()
}

func main() {
	// 判断文件名是否设置
	if path == "" {
		fmt.Println("-path or --path value is required")
	}
	// 获取文件属性
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Println("os.Stat error:", err)
	}

	fmt.Println("file name: ", fileInfo.Name())
	fmt.Println("file size: ", fileInfo.Size())

}
