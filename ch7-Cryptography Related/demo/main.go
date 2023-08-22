package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	// md5 sha1 sha256 sha512 使用方式都类似
	// 散列值一般是一个二进制的字符串,有些可见有些不可见,需要格式化
	// 格式化为16进制的数字串 0-9 a-f
	// 数据转换完成后,长度会是原来的2倍
	// 第一种方式
	data := []byte("These pretzels are making me thirsty.")
	fmt.Printf("%x", md5.Sum(data))
	fmt.Println()
	// 第二种方式
	h := md5.New()
	io.WriteString(h, "The fog is getting thicker!")
	io.WriteString(h, "And Leon's getting laaarger!")
	res := h.Sum(nil)
	fmt.Printf("%x", res)
	fmt.Println()
	fmt.Println(hex.EncodeToString(res))
	// 第三种方式
	f, err := os.Open("file.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	hx := md5.New()
	if _, err := io.Copy(hx, f); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%x", hx.Sum(nil))
	fmt.Println(hex.EncodeToString(res))
}
