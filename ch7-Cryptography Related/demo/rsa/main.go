package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"

	"github.com/767829413/advanced-go/util"
)

// rsa加密,公钥进行加密 /tmp/key/public.pem
func RSAEncrypt(plainText []byte, pubKeyFile string) ([]byte, error) {
	// 读取公钥文件内容
	buf, err := util.ReadFile(pubKeyFile)
	if err != nil {
		return nil, err
	}

	// pem解码
	block, _ := pem.Decode(buf)
	// x509规范解码
	pubAny, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pubKey, ok := pubAny.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key type conversion failed")
	}
	// 使用公钥加密
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, plainText)
	if err != nil {
		return nil, err
	}
	return cipherText, nil
}

// rsa解密,私钥进行解密 /tmp/key/private.pem
func RSADecrypt(cipherText []byte, priKeyFile string) ([]byte, error) {
	// 读取私钥文件内容
	buf, err := util.ReadFile(priKeyFile)
	if err != nil {
		return nil, err
	}

	// pem解码
	block, _ := pem.Decode(buf)
	// x509规范解码
	priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 使用私钥解密
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, priKey, cipherText)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}

var priKeyPath = "/tmp/key/private.pem"
var pubKeyPath = "/tmp/key/public.pem"

func GenerateHmacSHA512(plainText []byte) []byte {
	mHash := sha512.New()
	mHash.Write(plainText)
	hashText := mHash.Sum(nil)
	return hashText
}

// RSA签名
func SignatureRSA(plainText []byte, priKeyPath string) ([]byte, error) {
	/*
	   1. 从指定位置获取私钥文件
	   2. 输出私钥文件内容
	   3. 使用pem对私钥文件内容进行解码
	   4. 利用x509将数据还原为私钥内容
	   5. 创建一个哈希对象 -> md5|sha,计算散列值
	   6. 使用rsa相关函数(SignPKCS1v15)对散列值签名
	*/
	buf, err := util.ReadFile(priKeyPath)
	if err != nil {
		return nil, err
	}
	// pem解码
	block, _ := pem.Decode(buf)
	// x509规范解码
	priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	hashText := GenerateHmacSHA512(plainText)
	return rsa.SignPKCS1v15(rand.Reader, priKey, crypto.SHA512, hashText)
}

// RSA签名校验
func VerifySignatureRSA(plainText, sigText []byte, pubKeyPath string) (bool, error) {
	/*
		1. 从指定位置获取公钥文件
		2. 输出公钥文件内容
		3. 使用pem对私钥文件内容进行解码
		4. 利用x509将数据还原为公钥内容
		5. 将内容进行断言->得到公钥结构体
		6. 创建一个哈希对象(和签名过程的算法一致) -> md5|sha,计算散列值
		7. 使用rsa相关函数(VerifyPKCS1v15)对签名进行验证
	*/
	buf, err := util.ReadFile(pubKeyPath)
	if err != nil {
		return false, err
	}
	// pem解码
	block, _ := pem.Decode(buf)
	// x509规范解码
	pubAny, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, err
	}
	pubKey, ok := pubAny.(*rsa.PublicKey)
	if !ok {
		return false, fmt.Errorf("public key type conversion failed")
	}
	hashText := GenerateHmacSHA512(plainText)
	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA512, hashText, sigText)
	if err != nil {
		return false, err
	}
	return true, nil
}

func main() {
	// /tmp/key/private.pem
	// /tmp/key/public.pem
	util.GenerateRsaKey(1024, "/tmp/key/private.pem", "/tmp/key/public.pem")
	cipherText, err := RSAEncrypt([]byte("你是一只巨"), "/tmp/key/public.pem")
	if err != nil {
		fmt.Println(err)
		return
	}
	plainText1, err := RSADecrypt(cipherText, "/tmp/key/private.pem")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(plainText1))

	plainText := []byte("哈哈健康撒谎看哈上课回答是科技活动空间")
	sig, err := SignatureRSA(plainText, priKeyPath)
	if err != nil {
		log.Println(err)
		return
	}
	res, err := VerifySignatureRSA(plainText, sig, pubKeyPath)
	log.Println(res, err)

}
