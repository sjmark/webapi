package utils

import (
	"crypto/x509"
	"crypto/rsa"
	"encoding/pem"
	"errors"
	"crypto/rand"
	"encoding/base64"
)

// 创建私钥:
// openssl genrsa -out private.pem 1024
// 创建公钥:
// openssl rsa -in private.pem -pubout -out public.pem

// GOLANG语言也提供了API用于生成私钥
// encoding/pem
// crypto/x509

// PKCS(公钥密码标准)

// 公钥和私钥可以从文件中读取
var privateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDZsfv1qscqYdy4vY+P4e3cAtmvppXQcRvrF1cB4drkv0haU24Y
7m5qYtT52Kr539RdbKKdLAM6s20lWy7+5C0DgacdwYWd/7PeCELyEipZJL07Vro7
Ate8Bfjya+wltGK9+XNUIHiumUKULW4KDx21+1NLAUeJ6PeW+DAkmJWF6QIDAQAB
AoGBAJlNxenTQj6OfCl9FMR2jlMJjtMrtQT9InQEE7m3m7bLHeC+MCJOhmNVBjaM
ZpthDORdxIZ6oCuOf6Z2+Dl35lntGFh5J7S34UP2BWzF1IyyQfySCNexGNHKT1G1
XKQtHmtc2gWWthEg+S6ciIyw2IGrrP2Rke81vYHExPrexf0hAkEA9Izb0MiYsMCB
/jemLJB0Lb3Y/B8xjGjQFFBQT7bmwBVjvZWZVpnMnXi9sWGdgUpxsCuAIROXjZ40
IRZ2C9EouwJBAOPjPvV8Sgw4vaseOqlJvSq/C/pIFx6RVznDGlc8bRg7SgTPpjHG
4G+M3mVgpCX1a/EU1mB+fhiJ2LAZ/pTtY6sCQGaW9NwIWu3DRIVGCSMm0mYh/3X9
DAcwLSJoctiODQ1Fq9rreDE5QfpJnaJdJfsIJNtX1F+L3YceeBXtW0Ynz2MCQBI8
9KP274Is5FkWkUFNKnuKUK4WKOuEXEO+LpR+vIhs7k6WQ8nGDd4/mujoJBr5mkrw
DPwqA3N5TMNDQVGv8gMCQQCaKGJgWYgvo3/milFfImbp+m7/Y3vCptarldXrYQWO
AQjxwc71ZGBFDITYvdgJM1MTqc8xQek1FXn1vfpy2c6O
-----END RSA PRIVATE KEY-----
`)

var publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDZsfv1qscqYdy4vY+P4e3cAtmv
ppXQcRvrF1cB4drkv0haU24Y7m5qYtT52Kr539RdbKKdLAM6s20lWy7+5C0Dgacd
wYWd/7PeCELyEipZJL07Vro7Ate8Bfjya+wltGK9+XNUIHiumUKULW4KDx21+1NL
AUeJ6PeW+DAkmJWF6QIDAQAB
-----END PUBLIC KEY-----
`)

func RsaEncrypt(origData []byte) ([]byte, error) { // 加密

	block, _ := pem.Decode(publicKey)

	if block == nil {
		return nil, errors.New("public key error")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)

	if err != nil {
		return nil, err
	}

	pub := pubInterface.(*rsa.PublicKey)

	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)

}

func RsaDecrypt(ciphertext []byte) ([]byte, error) { // 解密

	block, _ := pem.Decode(privateKey)

	if block == nil {
		return nil, errors.New("private key error!")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	if err != nil {
		return nil, err
	}

	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)

}

func RsaBase64Encrypt(origData []byte) (string, error) { // 会将结果用BASE64进行编码

	result, err := RsaEncrypt(origData)

	if err != nil {
		return "", err
	}

	encode := base64.StdEncoding.EncodeToString(result)

	return encode, nil

}

func RsaBase64Decrypt(encode string) ([]byte, error) { // 对用BASE64编码的字符串进行解密

	ciphertext, err := base64.StdEncoding.DecodeString(encode)

	if err != nil {

		return nil, err

	}

	return RsaDecrypt(ciphertext)

}
func UnSafeRsaDecrypt(encode string) []byte { // 对用BASE64编码的字符串进行解密

	ciphertext, err := RsaBase64Decrypt(encode)

	if err != nil {

		return nil

	}

	return ciphertext

}
func RemoveZeroFromEnd(origin []byte) []byte {

	point := int(0)

	for idx, b := range origin {
		if b == 0 {
			point = idx
			break
		}
	}

	if point == 0 {
		return origin
	}

	return origin[:point]
}
