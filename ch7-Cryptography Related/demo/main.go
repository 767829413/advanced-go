package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func readFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fInfo, err := f.Stat()
	if err != nil {
		return nil, err
	}
	buf := make([]byte, fInfo.Size())
	f.Read(buf)
	return buf, nil
}

// 生成rsa的密钥对,保存到文件
func GenerateRsaKey(rsaKeyLen int) error {
	// 私钥生成流程
	priKey, err := rsa.GenerateKey(rand.Reader, rsaKeyLen)
	if err != nil {
		return err
	}
	derText := x509.MarshalPKCS1PrivateKey(priKey)
	blockPri := &pem.Block{
		Type:  "rsa private key",
		Bytes: derText,
	}
	// 创建文件流句柄
	fPri, err := os.Create("/tmp/key/private.pem")
	if err != nil {
		return err
	}
	defer fPri.Close()
	err = pem.Encode(fPri, blockPri)
	if err != nil {
		return err
	}
	// 公钥生成流程
	derStream, err := x509.MarshalPKIXPublicKey(&priKey.PublicKey)
	if err != nil {
		return err
	}
	blockPub := &pem.Block{
		Type:  "rsa public key",
		Bytes: derStream,
	}
	fPub, err := os.Create("/tmp/key/public.pem")
	if err != nil {
		return err
	}
	defer fPub.Close()
	err = pem.Encode(fPub, blockPub)
	if err != nil {
		return err
	}
	return nil
}

// rsa加密,公钥进行加密 /tmp/key/public.pem
func RSAEncrypt(plainText []byte, pubKeyFile string) ([]byte, error) {
	// 读取公钥文件内容
	buf, err := readFile(pubKeyFile)
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
	buf, err := readFile(priKeyFile)
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

func main() {
	GenerateRsaKey(4096)
	cipherText, err := RSAEncrypt([]byte("你是一只??????????asdsadsaasdsadsaasdsadsaasdsadsaasdsadsaasdsadsaasdsadsaasdsadsaasdsadsaasdsadsaasdsadsaasdsadsaasdsadsaasdsadsaasdsadsaasdsadsaasdsadsaasdsadsaasdsadsaasdsadsaasdsadsaasdsadsaasdsadsaassdsadsaasdsadsaasdsadsaasdsadsaasdsadsaasdsadsaasdsadsaasdsadsaasdsadsaasdsadsaasdsadsaasdsadsaasd"), "/tmp/key/public.pem")
	if err != nil {
		fmt.Println(err)
		return
	}
	plainText, err := RSADecrypt(cipherText, "/tmp/key/private.pem")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(plainText))
}
