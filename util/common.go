package util

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"hash/fnv"
	mrand "math/rand"
	"os"
	"path/filepath"
	"time"
	"unsafe"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var src = mrand.NewSource(time.Now().UnixNano())

const (
	// 6 bits to represent a letter index
	letterIdBits = 6
	// All 1-bits as many as letterIdBits
	letterIdMask = 1<<letterIdBits - 1
	letterIdMax  = 63 / letterIdBits
)

func RandStr(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(letters) {
			b[i] = letters[idx]
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

func Str2HashInt(s string) uint64 {
	h := fnv.New64()
	h.Write([]byte(s))
	return h.Sum64()
}

func ErrExit(err error, info string) {
	if err != nil {
		fmt.Println(info+":", err)
		os.Exit(1)
	}
}

func getExecPath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return exPath
}

func ReadFile(path string) ([]byte, error) {
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
// /tmp/key/private.pem
// /tmp/key/public.pem
func GenerateRsaKey(rsaKeyLen int, priPath, pubPath string) error {
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
	fPri, err := os.Create(priPath)
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
	fPub, err := os.Create(pubPath)
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

// 生成ECC密钥对
func GenerateEccKey(c elliptic.Curve, priPath, pubPath string) error {
	// 私钥生成流程
	// 使用crypto/ecdsa的(GenerateKey)来生成密钥对
	priKey, err := ecdsa.GenerateKey(c, rand.Reader)
	if err != nil {
		return err
	}
	//1. 使用x509(MarshalECPrivateKey)将私钥序列化
	derText, err := x509.MarshalECPrivateKey(priKey)
	if err != nil {
		return err
	}
	//2. 将序列化的数据放到pem.Block结构体中
	blockPri := &pem.Block{
		Type:  "ecdsa private key",
		Bytes: derText,
	}
	//3. 使用pem.Encode()编码
	fPri, err := os.Create(priPath)
	if err != nil {
		return err
	}
	defer fPri.Close()
	pem.Encode(fPri, blockPri)

	// 公钥生成流程
	// 1. 使用x509(MarshalPKIXPublicKey)将私钥序列化
	derStream, err := x509.MarshalPKIXPublicKey(&priKey.PublicKey)
	if err != nil {
		return err
	}
	// 2. 将序列化的数据放到pem.Block结构体中
	blockPub := &pem.Block{
		Type:  "ecdsa public key",
		Bytes: derStream,
	}
	// 3. 使用pem.Encode()编码
	fPub, err := os.Create(pubPath)
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
