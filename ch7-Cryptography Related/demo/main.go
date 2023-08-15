package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"fmt"
)

/*
DES的CBC加密
1. 编写填充函数,如果最后一个分组字节数不够,填充
2. 字节数合适的便添加新分组
3. 填充的字节值 == 减少的字节值
*/

func paddingLastGroup(plainText []byte, blockSize int) []byte {
	// 计算最后一组中剩余字节数,通过取余获取,恰好就填充整个一组
	padNum := blockSize - len(plainText)%blockSize
	// 创建新的byte切片,长度为panNum,每个字节值为byte(padNum)
	char := []byte{byte(padNum)}
	// 新的切片初始化
	char = bytes.Repeat(char, padNum)
	plainText = append(plainText, char...)
	return plainText
}

func unpaddingLastGroup(plainText []byte) []byte {
	// 获取最后一位获取填充长度
	l := int(plainText[len(plainText)-1])
	return plainText[:len(plainText)-l]
}

// des加密,分组方法CBC,key长度是8
func desEnCrypt(plainText, key []byte) ([]byte, error) {
	// 创建一个底层使用的 DES 的密码接口
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 根据分组模式进行分组填充(比如CBC,ECB需要填充)
	plainText = paddingLastGroup(plainText, block.BlockSize())
	// 创建一个密码分组模式的接口对象,这里是CBC
	iv := []byte("12345678")
	blockMode := cipher.NewCBCEncrypter(block, iv)
	dst := make([]byte, len(plainText))
	blockMode.CryptBlocks(dst, plainText)
	return dst, nil
}

// des解密,分组方法CBC,key长度是8
func desDecrypter(cipherText, key []byte) ([]byte, error) {
	// 创建一个底层使用的 DES 的密码接口
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 创建一个密码分组模式的接口对象,这里是CBC
	iv := []byte("12345678")
	blockMode := cipher.NewCBCDecrypter(block, iv)
	dst := make([]byte, len(cipherText))
	blockMode.CryptBlocks(dst, cipherText)
	return unpaddingLastGroup(dst), nil
}

// aes加密,分组方法CTR
func aesEnCrypt(plainText, key []byte) ([]byte, error) {
	// 创建一个底层使用的 AES 的密码接口
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 创建一个密码分组模式的接口对象,这里是CBC
	iv := []byte("1234567812345678")
	blockMode := cipher.NewCTR(block, iv)
	dst := make([]byte, len(plainText))
	blockMode.XORKeyStream(dst, plainText)
	return dst, nil
}

// aes解密,分组方法CTR
func aesDecrypter(cipherText, key []byte) ([]byte, error) {
	// 创建一个底层使用的 AES 的密码接口
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// 创建一个密码分组模式的接口对象,这里是CBC
	iv := []byte("1234567812345678")
	blockMode := cipher.NewCTR(block, iv)
	dst := make([]byte, len(cipherText))
	blockMode.XORKeyStream(dst, cipherText)
	return dst, nil
}

func main() {
	cipherText, _ := desEnCrypt([]byte("qwerweqrwertwe"), []byte("88888888"))
	plainText, _ := desDecrypter(cipherText, []byte("88888888"))
	fmt.Println(cipherText)
	fmt.Println(string(plainText) == "qwerweqrwertwe")
	cipherText, _ = aesEnCrypt([]byte("qwerweqrwertwe"), []byte("8888888888888888"))
	plainText, _ = aesDecrypter(cipherText, []byte("8888888888888888"))
	fmt.Println(cipherText)
	fmt.Println(string(plainText) == "qwerweqrwertwe")
}
