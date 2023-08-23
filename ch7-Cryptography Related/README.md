# 加密相关

1. 加密三要素
	* 明文/密文
	* 密钥
		* 定长字符串
		* 需要根据加密算法确定长度
	* 算法
		* 加密算法
		* 解密算法
		* 加密和解密算法可能互逆也可能相同
2. 常用的两种加密方式
	* 对称加密
		* 密钥: 加密,解密使用相同密钥
		* 特点
			* 数据的机密性能双方向保证
			* 加密效率高,适合大文件,大数据
			* 加密强度不高,相对于非对称加密
	* 非对称加密
		* 密钥: 加密,解密使用不同密钥,需要使用密钥生成算法来获取密钥对
			* 公钥: 可以公开的
				* 公钥加密需要对应私钥解密
			* 私钥: 需要进行妥善保管
				* 私钥加密,私钥解密
		* 特点
			* 数据的机密性只能单方向保证
			* 加密效率低,适合少量数据
			* 加密强度高,相对于对称加密

3. 密码安全常识
	* 不要使用保密的密码算法(普通公司或个人) 
	* 使用低强度密码比不进行任何加密更危险
	* 任何密码都有破解的一天
	* 密码只是信息安全中的一部分

## 对称加密
 
 **以分组为单位进行处理的密码算法称为分组密码**

1. 编码的概念
 
 1G = 1024m 1m = 1024kbyte 1byte = 8bit bit 0/1 
 (b byte B bit)
 
 **计算机的操作对象并不是文字,而是由0和1排列的比特序列,将现实世界中的东西映射为比特序列的操作称为编码**

 加密 => 编码 解密 => 解码

2. DES -- Data Encryption Standard
	* 什么是DES（Data Encryption Standard）:[资料加密标准(DES)](https://zh.wikipedia.org/zh-hans/è³æå å¯æ¨æº)

	* 加密和解密

	```text
	DES是一种将64比特的明加密成64比特的密文的对称密码算法，它的密钥长度是56比特。尽管从规格上来说，DES的密钥长度是64比特，但由于每隔7比特会设置一个用于错误检查的比特，因此实质上其密钥长度是56比特。

    DES是以64比特的明文(比特序列)为一个单位来进行加密的，这个64比特的单位称为分组。一般来说，以分组为单位进行处理的密码算法称为分组密码 (blockcipher)，DES就是分组码的一种。
    
	DES每次只能加密64比特的数据，如果要加密的明文比较长，就需要对DES加密进行迭代(反复)，而迭代的具体方式就称为模式(mode)。
	```

	![1.png](https://s2.loli.net/2023/08/11/AWesTMbROmJYjdK.png)

	* 使用DES方式加密安全吗?
		* 不安全,已经破解
	* 是不是分组密码?
		* 是,先对数据分组,然后加密解密
	* DES的分组长度?
		* 8byte = 64bit
	* DES的密钥长度?
		* 56bit密钥长度 + 8bit错误检测标志位 = 64bit = 8byte

3. 3DES -- TripleDES
	* 什么是3DES（Triple DES）: [三重数据加密算法（英语：Triple Data Encryption Algorithm，缩写为TDEA，Triple DEA）](https://zh.wikipedia.org/zh-hans/3DES)
	* 加密和解密
		* 加密 ![3DES-1.png](https://s2.loli.net/2023/08/11/aQsbDKJUfrkhwit.png)
		* 解密 ![3DES-2.png](https://s2.loli.net/2023/08/11/CHtJeBcAsjQZ1dh.png)
	* 使用3DES方式加密安全吗?
		* 安全,但是效率低
	* 是不是分组密码?
		* 是
	* 3DES的分组长度?
		* 8byte
	* 3DES的密钥长度?
		* 24byte,在算法内部会平均分成3份
	* 3DES的加密过程?
		* 密钥1加密,密钥2解密,密钥3加密
	* 3DES的解密过程?
		* * 密钥1解密,密钥2加密,密钥3解密

3. AES -- Advanced Encryption Standard
	* 什么是AES（Advanced Encryption Standard）: [高级加密标准（英语：Advanced Encryption Standard，缩写：AES），又称Rijndael加密法](https://zh.wikipedia.org/wiki/é«çº§å å¯æ å)
	* 使用AES方式加密安全吗?
		* 安全,效率高,推荐
	* 是不是分组密码?
		* 是
	* AES的分组长度?
		* 16byte = 128bit
	* AES的密钥长度?
		* 16byte = 128bit
		* 24byte = 192bit
		* 32byte = 256bit
		* go目前使用的是16byte

4. 分组密码模式
	* 维基百科: [分组密码工作模式](https://zh.wikipedia.org/wiki/åç»å¯ç å·¥ä½æ¨¡å¼)
	* 按位异或
		* 数据转换为二进制
		* 按位异或的操作符: ^
		* 两个标志位进行按位异或
			* 相同为0,不同为1
	* ECB- Electronic Code Book,电子密码本模式
		* 特点: 简单,高效,密文有规律,易破解
		* 最后一个明文分组必须填充
			* des/3des: 最后一个分组填充满 8byte
			* aes: 最后一个分组填充满 16byte
		* 不需要初始化向量
	* CBC- Cipher Block Chaining,密码块链模式
		* 特点: 密文无规律,使用率高
		* 最后一个明文分组必须填充
			* des/3des: 最后一个分组填充满 8byte
			* aes: 最后一个分组填充满 16byte
		* 需要初始化向量(数组)
			* 数组长度: 明文分组长度相同
			* 数据来源: 负责加密方提供(随机字符串)
			* 解密和加密的初始化向量必须相同
	* CFB- Cipher FeedBack,密文反馈模式
		* 特点: 密文无规律,明文分组是和一个数据流进行按位异或操作后最终生成密文
		* 最后一个明文分组不必填充
		* 需要初始化向量(数组)
			* 数组长度: 明文分组长度相同
			* 数据来源: 负责加密方提供(随机字符串)
			* 解密和加密的初始化向量必须相同
	* OFB - Output-Feedback,输出反馈模式
		* 特点: 密文无规律,明文分组是和一个数据流进行按位异或操作后最终生成密文
		* 最后一个明文分组不必填充
		* 需要初始化向量(数组)
			* 数组长度: 明文分组长度相同
			* 数据来源: 负责加密方提供(随机字符串)
			* 解密和加密的初始化向量必须相同
	* CTR-CounTeR,计数器模式
		* 特点: 密文无规律,明文分组是和一个数据流进行按位异或操作后最终生成密文
		* 最后一个明文分组不必填充
		* 不需要初始化向量
			* go接口中的IV可以理解为随机数种子,长度是明文分组长度
	* 最后一个明文分组的填充
		* 使用CBC,ECB分组模式需要填充
			* 要求: 
				* 明文分组中进行填充,然后加密
				* 解密密文得到明文,需要删除填充字节
				* 小技巧,填充的字节最好就是填充的长度值,如果明文分组不需要填充,那么也填充一个分组,方便删除
		* 使用OFB,CFB,CTR不需要填充
	* 初始化向量-IV
		* ECB,CTR分组模式不需要初始化向量
		* CBC,OFC,CFB需要初始化向量
			* 初始化向量长度
				* DES/3DES: 8byte
				* AES: 16byte
			* 加密解密的初始化向量是一致的

5. 对称加密在go中的实现

	* 加密流程
		1. 创建一个底层使用的 DES/3DES/AES的密码接口
			* [DES/3DES](https://pkg.go.dev/crypto/des@go1.21.0)
			* [AES](https://pkg.go.dev/crypto/aes@go1.21.0)
		2. 根据分组模式进行分组填充(比如CBC,ECB需要填充)
		3. 创建一个密码分组模式的接口对象
			* [CBC|CFB|OFB|CTR](https://pkg.go.dev/crypto/cipher@go1.21.0#Block)
		4. 加密得到密文

		```go
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
		```

## 非对称加密

1. 对称加密的弊端
	* 密钥分发困难
	* 通过非对称加密完成密钥分发

2. 非对称加密的密钥
	* 不存在密钥分发困难问题
	* 场景分析
		* 信息加密(A写数据给B,只允许B读)
			* A: 公钥 B: 私钥
		* 登陆认证(客户端登陆,请求服务器,向服务器请求个人数据)
			* 服务器: 公钥 客户端: 私钥
		* 数字签名(表明信息的真实性,附在信息原文后)
			* 发送信息的人: 私钥 收到信息的人: 公钥
		* 网银U盾 
			* 个人: 私钥 银行: 公钥
		* 总结: 数据对谁更重要,谁拿私钥
		* 直观上私钥比公钥长,一般生成的文件xxx.pub 公钥 xxx 私钥

3. 使用RSA非对称加密通信流程

```lua
            +---------+                    +---------+
            | Sender  |                    | Receiver|
            +---------+                    +---------+
                |                                |
                |           生成密钥对            |
                +------------------------------> |
                |                                |
                |          请求公钥               |
                +------------------------------> |
                |                                |
                |          返回公钥               |
                | <------------------------------+
                |                                |
                |        加密数据                  |
                | -----------------------------> |
                |                                |
                |        使用公钥加密数据          |
                | -----------------------------> |
                |                                |
                |        使用私钥解密数据          |
                | <----------------------------+ |
                |                                |
                |          返回解密后的数据        |
                | <----------------------------+ |
```

4. 生成RSA的密钥对

	* [RSA加密算法](https://zh.wikipedia.org/wiki/RSAå å¯æ¼ç®æ³)
	* [Golang中RSA相关package](https://pkg.go.dev/crypto/rsa)
	* [Golang中x509相关package](https://pkg.go.dev/crypto/x509)
	* [Golang中pem相关package](https://pkg.go.dev/encoding/pem)
	* 生成私钥操作流程
		* 使用crypto中的rsa相关的方法生成私钥
			* func GenerateKey(random io.Reader, bits int) (priv *PrivateKey, err error)
			* rand.Reader
			* 生成位数建议为1024整数倍
		* 通过x509标准将得到的rsa私钥序列化为ASN.1的DER编码字符串
			* func MarshalPKCS1PrivateKey(key *rsa.PrivateKey) []byte
		* 将私钥字符串设置到pem格式块中
			* 初始化一个pem.Block结构体

			```go
			type Block struct {
				Type    string            // The type, taken from the preamble (i.e. "RSA PRIVATE KEY").
				Headers map[string]string // Optional headers.
				Bytes   []byte            // The decoded bytes of the contents. Typically a DER encoded ASN.1 structure.
			}
			```

		* 通过pem将设置好的数据进行编码,并写入磁盘文件
			* func Encode(out io.Writer, b *Block) error
			* out: 指定一个文件指针就行
	* 生成公钥流程
		* 从得到的私钥对象中将公钥信息提取

		```go
		type PrivateKey struct {
			PublicKey            // public part.
			D         *big.Int   // private exponent
			Primes    []*big.Int // prime factors of N, has >= 2 elements.

			// Precomputed contains precomputed values that speed up RSA operations,
			// if available. It must be generated by calling PrivateKey.Precompute and
			// must not be modified.
			Precomputed PrecomputedValues
		}

		type PublicKey struct {
			N *big.Int // modulus
			E int      // public exponent
		}
		```

		* 通过x509标准将得到的rsa公钥序列化为字符串
			* func MarshalPKIXPublicKey(pub any) ([]byte, error)
		* 将公钥字符串设置到pem格式块中
		* 通过pem将设置好的数据进行编码,并写入磁盘文件

		```go
		package main

		import (
			"crypto/rand"
			"crypto/rsa"
			"crypto/x509"
			"encoding/pem"
			"os"
		)

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

		func main() {
			GenerateRsaKey(1024)
		}
		```

5. RSA加解密
	* 加密
		* 将公钥文件中的公钥读出,得到使用pem编码的字符串
			* 读文件
		* 将得到的字符串解码
			* pem.Decode
		* 使用x509将编码后的公钥解析出来
			* func ParsePKCS1PublicKey(der []byte) (*rsa.PublicKey, error)
		* 使用得到的公钥通过rsa进行加密
			* func EncryptPKCS1v15(random io.Reader, pub *PublicKey, msg []byte) ([]byte, error)
	* 解密

	```go
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
		cipherText, err := RSAEncrypt([]byte("你是一只巨"), "/tmp/key/public.pem")
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
	```

6. ECC 椭圆曲线
	* [椭圆曲线密码学](https://zh.wikipedia.org/wiki/%E6%A4%AD%E5%9C%86%E6%9B%B2%E7%BA%BF%E5%AF%86%E7%A0%81%E5%AD%A6s)
	* go中的椭圆曲线相关package
		* [crypto/elliptic](https://pkg.go.dev/crypto/elliptic@go1.21.0)
		* [crypto/ecdsa](https://pkg.go.dev/crypto/ecdsa@go1.21.0)

7. 对称加密和非对称加密比较
	* 对称加密
		* 使用相同的密钥进行加密和解密。
		* 加密和解密速度较快，适用于大量数据的加密。
		* 密钥的管理较为复杂，需要确保密钥的安全性。
		* 不适用于需要安全地传输密钥的场景。
	* 非对称加密
		* 使用两个密钥，一个用于加密，另一个用于解密。
		* 加密和解密速度较慢。
		* 密钥的管理相对简单，只需保护私钥的安全性。
		* 适用于需要安全地传输密钥的场景，例如数字证书和安全通信。
	* 总结：
		* 对称加密适用于需要高效加密和解密大量数据的场景
		* 非对称加密适用于需要安全地传输密钥和进行安全通信的场景。
		* 在实际应用中，通常会结合使用对称加密和非对称加密，以兼顾效率和安全性。

## 哈希函数

1. 概念

 单向散列函数,哈希函数,散列函数,消息摘要函数,杂凑函数都是一种

 [散列函數](https://zh.wikipedia.org/wiki/æ£åå½æ¸)

 接收的输入: 原像
 输出: 散列值,哈希值,指纹,摘要

2. 特性

	* 将任意长度的数据转换成固定成都的数据
	* 很强的抗碰撞性
	* 不可逆

3. 常用哈希函数

	* [MD4,MD5](https://zh.wikipedia.org/wiki/MD5)
		* 不安全
		* 散列值长度: 128bit == 16byte
	* [SHA家族](https://zh.wikipedia.org/wiki/SHA%E5%AE%B6%E6%97%8F)
		* 安全性
			* sha-1: 不安全
			* sha-2及以上: 目前安全
		* 散列值长度: 
			* sha-1: 160bit == 20byte
			* sha-224: 224bit == 28byte
			* sha-256: 256bit == 32byte
			* sha-384: 384bit == 48byte
			* sha-512: 512bit == 64byte

4. go中使用哈希函数

	* [crypto/md5](https://pkg.go.dev/crypto/md5@go1.21.0)
	* [crypto/sha1](https://pkg.go.dev/crypto/sha1@go1.21.0)
	* [crypto/sha256](https://pkg.go.dev/crypto/sha256@go1.21.0)
	* [crypto/sha512](https://pkg.go.dev/crypto/sha512@go1.21.0)

	```go
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

		h := md5.New()
		if _, err := io.Copy(h, f); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%x", h.Sum(nil))	
	}
	```

5. 消息认证码
	* 前提条件: 
		* 在消息认证码生成的一方和校验的一方必须有一个共同持有的密钥
		* 双方必须约定好使用同样的哈希函数进行数据计算
	* 流程: 
		* 发送者:
			* 发送原始消息
			* 将原始消息生成消息验证码
				* {{原始消息} + 密钥} * 哈希函数 = 散列值(消息认证码)
			* 将消息发送给对方
		* 接收者:
			* 接收原始数据
			* 接收消息认证码
			* 校验:
				* {{接收消息} + 密钥} * 哈希函数 = 新散列值
				* 校验新的散列值和接收的散列值是否一致
	* go中如何使用消息认证码
		* [crypto/hmac](https://pkg.go.dev/crypto/hmac)

		```go
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
		```

	* 消息认证码的缺点
		* 可解决问题
			* 密钥分发困难
		* 无法解决问题
			* 不能进行第三方证明
			* 不能防止否认

6. 数字签名
	* 过程:
		* 签名
			* 有原始数据对其进行哈希运算 -> 散列值
			* 使用非对称加密的私钥对散列值加密 -> 签名
			* 将原始数据和签名一并发送给对方
		* 签名验证
			* 接收数据
				* 原始数据
				* 数字签名
			* 数字签名,需要使用公钥进行解密,得到散列值
			* 对原始数据进行哈希运算得到新的散列值
			* 校验原始散列值和新的散列值是否一致
	* 非对称加密和数字签名
		* 总结: 
			1. 数据通信
				* 公钥加密,私钥解密
			2. 数字签名
				* 私钥加密,公钥解密
	* 使用RSA进行数字签名
		* 使用RSA生成密钥对
			* 生成密钥对
			* 序列化
			* 保存到相应位置持久化
		* 使用私钥进行数字签名
			1. 从指定位置获取私钥文件
			2. 输出私钥文件内容
			3. 使用pem对私钥文件内容进行解码
			4. 利用x509将数据还原为私钥内容
			5. 创建一个哈希对象 -> md5|sha,计算散列值
			6. 使用rsa相关函数(SignPKCS1v15)对散列值签名
		* 使用公钥进行认证
			1. 从指定位置获取公钥文件
			2. 输出公钥文件内容
			3. 使用pem对私钥文件内容进行解码
			4. 利用x509将数据还原为公钥内容
			5. 将内容进行断言->得到公钥结构体
			6. 创建一个哈希对象(和签名过程的算法一致) -> md5|sha,计算散列值
			7. 使用rsa相关函数(VerifyPKCS1v15)对签名进行验证
		* go中使用RSA进行数字签名

		```go
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
			// util.GenerateRsaKey(1024)
			plainText := []byte("哈哈健康撒谎看哈上课回答是科技活动空间")
			sig, err := SignatureRSA(plainText, priKeyPath)
			if err != nil {
				log.Println(err)
				return
			}
			res, err := VerifySignatureRSA(plainText, sig, pubKeyPath)
			log.Println(res, err)
		}
		```

	* 使用椭圆曲线进行数字签名
		* [椭圆曲线密码学](https://zh.wikipedia.org/wiki/%E6%A4%AD%E5%9C%86%E6%9B%B2%E7%BA%BF%E5%AF%86%E7%A0%81%E5%AD%A6)
		* go中的椭圆曲线相关package
			* [crypto/elliptic](https://pkg.go.dev/crypto/elliptic@go1.21.0)
			* [crypto/ecdsa](https://pkg.go.dev/crypto/ecdsa@go1.21.0)
		* 推荐使用的5个素域上的椭圆曲线素数模
			* P-192 = $2^{192}$ - $2^{64}$ - 1
			* P-224 = $2^{224}$ - $2^{96}$ + 1
			* P-256 = $2^{256}$ - $2^{224}$ + $2^{192}$ - $2^{96}$ - 1
			* P-384 = $2^{384}$ - $2^{128}$ - $2^{96}$ + $2^{32}$ - 1
			* P-521 = $2^{521}$ - 1
		* ECC密钥对生成流程
			1. 使用crypto/ecdsa的(GenerateKey)来生成密钥对
			2. 将私钥写入磁盘
				1. 使用x509(MarshalECPrivateKey)将私钥序列化
				2. 将序列化的数据放到pem.Block结构体中
				3. 使用pem.Encode()编码
			2. 将公钥写入磁盘
				1. 使用x509(MarshalPKIXPublicKey)将私钥序列化
				2. 将序列化的数据放到pem.Block结构体中
				3. 使用pem.Encode()编码
			3. 使用私钥进行数字签名
				1. 从指定位置获取私钥文件
				2. 输出私钥文件内容
				3. 使用pem对私钥文件内容进行解码
				4. 利用x509(ParseECPrivateKey)将数据还原为私钥内容
				5. 创建一个哈希对象 -> md5|sha,计算散列值
				6. 获取散列值
					* ecdsa相关函数(Sign)
						* func Sign(rand io.Reader, priv *PrivateKey, hash []byte) (r, s *big.Int, err error)
						* 使用ecdsa相关函数(Sign)对散列值签名,得到 r, s *big.Int
						* 得到的 r, s *big.Int 不能直接使用,需要进行序列化
						* 使用math相关函数(MarshalText)进行序列化
							* func (x *Int) MarshalText() (text []byte, err error)
					* ecdsa相关函数(SignASN1)
						* func SignASN1(rand io.Reader, priv *PrivateKey, hash []byte) ([]byte, error)
						* 得到散列值
			4. 使用公钥校验数字签名
				1. 从指定位置获取公钥文件
				2. 输出公钥文件内容
				3. 使用pem对私钥文件内容进行解码
				4. 利用x509(ParsePKIXPublicKey)将数据还原为公钥内容
				5. 将内容进行断言->得到公钥结构体

				```go
				type PublicKey struct {
					elliptic.Curve
					X, Y *big.Int
				}
				```

				6. 创建一个哈希对象(和签名过程的算法一致) -> md5|sha,计算散列值
				7. 进行验证
					* ecdsa相关函数(Verify)
						* func Verify(pub *PublicKey, hash []byte, r, s *big.Int) bool
						* 使用math相关函数(UnmarshalText)进行序列化
							* func (z *Int) UnmarshalText(text []byte) error
						* 获取验证结果
					* ecdsa相关函数(VerifyASN1)
						* func VerifyASN1(pub *PublicKey, hash, sig []byte) bool
						* 获取验证结果

			5. 在golang中代码实现

			```go
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
			```

			6. 数字签名的缺点
				* 无法确定公钥是否为消息发送方
					* 引入第三方认证机构(CA)

## 证书相关

1. 证书

2. SSL/TLS

3. https

4. 自签名证书
