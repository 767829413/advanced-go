package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"

	"github.com/767829413/advanced-go/util"
)

var priKeyPath = "/tmp/key/ecdsa/pri.pem"
var pubKeyPath = "/tmp/key/ecdsa/pub.pem"

func GenerateHmacSHA512(plainText []byte) []byte {
	mHash := sha512.New()
	mHash.Write(plainText)
	hashText := mHash.Sum(nil)
	return hashText
}

func signatureECCPre(plainText []byte, priKeyPath string) (hashText []byte, priKey *ecdsa.PrivateKey, err error) {
	//1. 从指定位置获取私钥文件
	//2. 输出私钥文件内容
	//3. 使用pem对私钥文件内容进行解码
	//4. 利用x509(ParseECPrivateKey)将数据还原为私钥内容
	//5. 创建一个哈希对象 -> md5|sha,计算散列值
	buf, err := util.ReadFile(priKeyPath)
	if err != nil {
		return nil, nil, err
	}
	// pem解码
	block, _ := pem.Decode(buf)
	// x509规范解码
	priKey, err = x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, nil, err
	}
	hashText = GenerateHmacSHA512(plainText)
	return hashText, priKey, nil
}

// ECC签名 - 私钥
func SignatureECCSign(plainText []byte, priKeyPath string) (r, s []byte, err error) {
	//* ecdsa相关函数(Sign)
	//* func Sign(rand io.Reader, priv *PrivateKey, hash []byte) (r, s *big.Int, err error)
	//* 使用ecdsa相关函数(Sign)对散列值签名,得到 r, s *big.Int
	//* 得到的 r, s *big.Int 不能直接使用,需要进行序列化
	//* 使用math相关函数(MarshalText)进行序列化
	//	* func (x *Int) MarshalText() (text []byte, err error)
	hashText, priKey, err := signatureECCPre(plainText, priKeyPath)
	if err != nil {
		return nil, nil, err
	}
	rBig, sBig, err := ecdsa.Sign(rand.Reader, priKey, hashText)
	if err != nil {
		return nil, nil, err
	}
	r, err = rBig.MarshalText()
	if err != nil {
		return nil, nil, err
	}
	s, err = sBig.MarshalText()
	if err != nil {
		return nil, nil, err
	}
	return r, s, nil
}

func SignatureECCSignASN1(plainText []byte, priKeyPath string) (sig []byte, err error) {
	//ecdsa相关函数(SignASN1)
	//* func SignASN1(rand io.Reader, priv *PrivateKey, hash []byte) ([]byte, error)
	//* 得到散列值
	hashText, priKey, err := signatureECCPre(plainText, priKeyPath)
	if err != nil {
		return nil, err
	}
	return ecdsa.SignASN1(rand.Reader, priKey, hashText)
}

func verifySignatureECCPre(plainText []byte, pubKeyPath string) ([]byte, *ecdsa.PublicKey, error) {
	//1. 从指定位置获取公钥文件
	//2. 输出公钥文件内容
	//3. 使用pem对私钥文件内容进行解码
	//4. 利用x509(ParsePKIXPublicKey)将数据还原为公钥内容
	//5. 将内容进行断言->得到公钥结构体

	//```go
	//type PublicKey struct {
	//	elliptic.Curve
	//	X, Y *big.Int
	//}
	//```

	//6. 创建一个哈希对象(和签名过程的算法一致) -> md5|sha,计算散列值
	buf, err := util.ReadFile(pubKeyPath)
	if err != nil {
		return nil, nil, err
	}
	// pem解码
	block, _ := pem.Decode(buf)
	// x509规范解码
	pubAny, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, nil, err
	}
	pubKey, ok := pubAny.(*ecdsa.PublicKey)
	if !ok {
		return nil, nil, fmt.Errorf("public key type conversion failed")
	}
	hashText := GenerateHmacSHA512(plainText)
	return hashText, pubKey, nil
}

// ECC校验签名 - 公钥
func VerifySignatureECCASN1(plainText, sig []byte, pubKeyPath string) (bool, error) {
	//* ecdsa相关函数(VerifyASN1)
	//* func VerifyASN1(pub *PublicKey, hash, sig []byte) bool
	//* 获取验证结果
	hashText, pubKey, err := verifySignatureECCPre(plainText, pubKeyPath)
	if err != nil {
		return false, nil
	}
	return ecdsa.VerifyASN1(pubKey, hashText, sig), nil

}

func VerifySignatureECC(plainText, r, s []byte, pubKeyPath string) (bool, error) {
	//* ecdsa相关函数(Verify)
	//* func Verify(pub *PublicKey, hash []byte, r, s *big.Int) bool
	//* 使用math相关函数(UnmarshalText)进行序列化
	//	* func (z *Int) UnmarshalText(text []byte) error
	//* 获取验证结果
	hashText, pubKey, err := verifySignatureECCPre(plainText, pubKeyPath)
	if err != nil {
		return false, nil
	}
	var rBig, sBig big.Int
	err = rBig.UnmarshalText(r)
	if err != nil {
		return false, nil
	}
	err = sBig.UnmarshalText(s)
	if err != nil {
		return false, nil
	}
	return ecdsa.Verify(pubKey, hashText, &rBig, &sBig), nil
}

func main() {
	util.GenerateEccKey(elliptic.P384(), "/tmp/key/ecdsa/pri.pem", "/tmp/key/ecdsa/pub.pem")

	plainText1 := []byte("sdsdsdsdsdssds")
	plainText2 := []byte("gffgfdhfgjfgjj")

	r, s, err := SignatureECCSign(plainText1, priKeyPath)
	if err != nil {
		fmt.Println("SignatureECCSign error: ", err)
		return
	}
	fmt.Println("VerifySignatureECC:")
	fmt.Println(VerifySignatureECC(plainText1, r, s, pubKeyPath))
	fmt.Println()
	sig, err := SignatureECCSignASN1(plainText2, priKeyPath)
	if err != nil {
		fmt.Println("SignatureECCSignASN1 error: ", err)
		return
	}
	fmt.Println("VerifySignatureECCASN1:")
	fmt.Println(VerifySignatureECCASN1(plainText2, sig, pubKeyPath))
}
