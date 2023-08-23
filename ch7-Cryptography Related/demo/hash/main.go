package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	// md5 sha1 sha256 sha512 使用方式都类似
	// 第一种方式
	data := []byte("These pretzels are making me thirsty.")
	fmt.Printf("%x", md5.Sum(data))
	// 第二种方式
	h := md5.New()
	io.WriteString(h, "The fog is getting thicker!")
	io.WriteString(h, "And Leon's getting laaarger!")
	fmt.Printf("%x", h.Sum(nil))
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
}
