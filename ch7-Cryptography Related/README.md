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
	* 什么是DES（Data Encryption Standard）:[资料加密标准(DES)](https://zh.wikipedia.org/zh-hans/%E8%B3%87%E6%96%99%E5%8A%A0%E5%AF%86%E6%A8%99%E6%BA%96)

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
	* 什么是AES（Advanced Encryption Standard）: [高级加密标准（英语：Advanced Encryption Standard，缩写：AES），又称Rijndael加密法](https://zh.wikipedia.org/wiki/%E9%AB%98%E7%BA%A7%E5%8A%A0%E5%AF%86%E6%A0%87%E5%87%86)
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
	* 维基百科: [分组密码工作模式](https://zh.wikipedia.org/wiki/%E5%88%86%E7%BB%84%E5%AF%86%E7%A0%81%E5%B7%A5%E4%BD%9C%E6%A8%A1%E5%BC%8F)
	* 按位异或
		* 数据转换为二进制
		* 按位异或的操作符: ^
		* 两个标志位进行按位异或
			* 相同为0,不同为1
	* ECB- Electronic Code Book,电子密码本模式
	* CBC- Cipher Block Chaining,密码块链模式

	```go
	
	```

	* CFB- Cipher FeedBack,密文反馈模式
	* OFB - Output-Feedback,输出反馈模式
	* CTR-CounTeR,计数器模式
	* 最后一个明文分组的填充
	* 初始化向量-IV

## 非对称加密
