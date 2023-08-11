package main

import (
	"bytes"
)

/*
DES的CBC加密
1. 编写填充函数,如果最后一个分组字节数不够,填充
2. 字节数合适的便添加新分组
3. 填充的字节值 == 减少的字节值
*/

func paddingLastGroup(plainText []byte, blockSize int) []byte {
	// 计算最后一组中剩余字节数,通过取余获取
	padNum := blockSize - len(plainText)%blockSize
	// 创建新的byte切片,长度为panNum,每个字节值为byte(padNum)
	char := []byte{byte(padNum)}
	// 新的切片初始化
	bytes.Repeat(char, padNum)
	plainText = append(plainText, char...)
	return plainText
}

func main() {
	paddingLastGroup([]byte("dsadasdsa"),8)
}
