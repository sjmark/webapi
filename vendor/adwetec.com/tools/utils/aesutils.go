package utils

import (
	"bytes"
	"crypto/aes"
	"encoding/hex"
	"crypto/cipher"
)

// AES加密数据块分组长度必须为128BIT.密钥长度可以是128BIT、192BIT、256BIT中的任意一个
// ******************************************************************************
func Pkcs5Padding(ciphertext []byte, blockSize int) []byte {

	padding := blockSize - len(ciphertext)%blockSize

	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(ciphertext, padtext...)
}
func Pkcs5UnPadding(origData []byte) []byte {

	length := len(origData)

	unpadding := int(origData[length-1])

	return origData[:(length - unpadding)]
}
func ZeroPadding(ciphertext []byte, blockSize int) []byte {

	padding := blockSize - len(ciphertext)%blockSize

	padtext := bytes.Repeat([]byte{0}, padding)

	return append(ciphertext, padtext...)
}
func ZeroUnPadding(origData []byte) []byte {

	length := len(origData)

	unpadding := int(origData[length-1])

	return origData[:(length - unpadding)]
}
func AesEncrypt(origData, key []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()

	origData = Pkcs5Padding(origData, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])

	crypted := make([]byte, len(origData))

	blockMode.CryptBlocks(crypted, origData)

	return crypted, nil
}

func AesDecrypt(crypted, key []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()

	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])

	origData := make([]byte, len(crypted))

	blockMode.CryptBlocks(origData, crypted)

	origData = Pkcs5UnPadding(origData)

	return origData, nil
}

// ******************************************************************************
// 上面js代码最后返回的是16进制
// 所以收到的数据hexText还需要用hex.DecodeString(hexText)转一下，这里略了
func Decrypt(ciphertext, key []byte) ([]byte, error) {

	pkey := PaddingLeft(key, '0', 16)

	block, err := aes.NewCipher(pkey) //选择加密算法
	if err != nil {
		return nil, err
	}

	blockModel := cipher.NewCBCDecrypter(block, pkey)

	plantText := make([]byte, len(ciphertext))

	blockModel.CryptBlocks(plantText, []byte(ciphertext))

	plantText = PKCS7UnPadding(plantText, block.BlockSize())

	return plantText, nil
}

func PKCS7UnPadding(plantText []byte, blockSize int) []byte {

	length := len(plantText)

	unpadding := int(plantText[length-1])

	return plantText[:(length - unpadding)]
}

func PaddingLeft(ori []byte, pad byte, length int) []byte {

	if len(ori) >= length {
		return ori[:length]
	}

	pads := bytes.Repeat([]byte{pad}, length-len(ori))

	return append(pads, ori...)
}

// ******************************************************************************
func AesCtrEncrypt(origData, key []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCTR(block, key)

	crypted := make([]byte, len(origData))

	blockMode.XORKeyStream(crypted, origData)

	return crypted, nil
}

func AesCtrDecrypt(crypted, key []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCTR(block, key)

	origData := make([]byte, len(crypted))

	blockMode.XORKeyStream(origData, crypted)

	return origData, nil
}

func AesCtrEncryptString(origData, key []byte) (string, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockMode := cipher.NewCTR(block, key)

	crypted := make([]byte, len(origData))

	blockMode.XORKeyStream(crypted, origData)

	return hex.EncodeToString(crypted), nil
}

func AesCtrDecryptString(crypted string, key []byte) (string, error) {

	cryptbytes, err := hex.DecodeString(crypted)

	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	blockMode := cipher.NewCTR(block, key)

	origData := make([]byte, len(crypted))

	blockMode.XORKeyStream(origData, cryptbytes)

	return string(ZeroUnPadding(origData)), nil
}
