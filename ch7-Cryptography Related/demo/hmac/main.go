package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
)

func GenerateHmac(plainText, key []byte) []byte {
	h := hmac.New(sha1.New, key)
	h.Write(plainText)
	hashText := h.Sum(nil)
	return hashText
}

func VerifyHmac(plainText, key, hashText []byte) bool {
	h := hmac.New(sha1.New, key)
	h.Write(plainText)
	hashNewText := h.Sum(nil)
	return hmac.Equal(hashNewText, hashText)
}

func main() {
	plainText, key := []byte("哈哈哈哈哈哈哈"), []byte("123456")
	hashText := GenerateHmac(plainText, key)
	fmt.Println(VerifyHmac(plainText, key, hashText))
}
