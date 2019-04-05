package utils

import (
	"crypto/des"
	"crypto/cipher"
	"bytes"
)

// DES:
// 其入口参数有三个:
// key:  加密解密使用的密钥
// data: 加密解密的数据
// mode: 其工作模式
// 		-- 当模式为加密模式时明文按照64位进行分组形成明文组(KEY用于对数据加密)
//      -- 当模式为解密模式时KEY用于对数据解密

// 实际运用中密钥只用到了64位中的56位.这样才具有高的安全性.DES的常见变体是三重DES.使用168位的密钥对资料进行三次加密的一种机制

// 涉及到加密模式和填充方式

// 关于分组加密:分组密码每次加密一个数据分组(这个分组的位数可以是随意的.一般选择64或者128位).另一方面流加密程序每次可以加密或解密一个字节的数据这就使它比流加密的应用程序更为有用

// 在用DES加密解密时经常会涉及到一个概念:
// 1.块(BLOCK--也叫分组)
// 2.模式(比如CBC)
// 3.初始向量(IV)
// 4.填充方式(PADDING--包括NONE用'\0'填充\pkcs5padding或pkcs7padding)

// 采用3DES、CBC模式、PKCS5PADDING.初始向量用KEY充当.另外对于ZEROPADDING还得约定好.对于数据长度刚好是BLOCKSIZE的整数倍时是否需要额外填充

// 跨语言加密解密应该使用PKCS5PADDING填充

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {

	padding := blockSize - len(ciphertext)%blockSize

	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(ciphertext, padtext...)
}
func PKCS5UnPadding(origData []byte) []byte {

	length := len(origData)

	unpadding := int(origData[length-1])

	return origData[:(length - unpadding)]
}

func DesEncrypt(origData, key []byte) ([]byte, error) {

	block, err := des.NewCipher(key) // 表示加密
	if err != nil {
		return nil, err
	}

	origData = Pkcs5Padding(origData, block.BlockSize()) // 填充方式 origData = ZeroPadding(origData, block.BlockSize())

	blockMode := cipher.NewCBCEncrypter(block, key) // 加密模式

	crypted := make([]byte, len(origData))

	blockMode.CryptBlocks(crypted, origData)

	return crypted, nil
}

func DesDecrypt(crypted, key []byte) ([]byte, error) {

	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, key)

	origData := make([]byte, len(crypted)) // origData := crypted

	blockMode.CryptBlocks(origData, crypted)

	origData = PKCS5UnPadding(origData) // origData = ZeroUnPadding(origData)

	return origData, nil
}

// 3DES加密
func TripleDesEncrypt(origData, key []byte) ([]byte, error) { // 秘钥长度必须是24BYTE(不同语言会有所不同)

	// 初始化向量长度必须是8字节

	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}

	origData = PKCS5Padding(origData, block.BlockSize()) // origData = ZeroPadding(origData, block.BlockSize())

	blockMode := cipher.NewCBCEncrypter(block, key[:8])

	crypted := make([]byte, len(origData))

	blockMode.CryptBlocks(crypted, origData)

	return crypted, nil
}

// 3DES解密
func TripleDesDecrypt(crypted, key []byte) ([]byte, error) {

	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, key[:8])

	origData := make([]byte, len(crypted))

	blockMode.CryptBlocks(origData, crypted)

	origData = PKCS5UnPadding(origData) // origData = ZeroUnPadding(origData)

	return origData, nil
}
