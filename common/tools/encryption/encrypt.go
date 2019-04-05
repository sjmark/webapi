package encryption

import (
	"bytes"
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rc4"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	// "toast/rsa/privatekey"
)

func Rc4(p string, key string) string {
	k := []byte(key)
	srctmp := []byte(p)
	cl, _ := rc4.NewCipher(k)
	dst := make([]byte, len(srctmp))
	cl.XORKeyStream(dst, srctmp)
	str := string(dst)
	return str
}

func Base64(s string) string {
	str, _ := base64.StdEncoding.DecodeString(s)
	return string(str)
}

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

//RSA加密
func RsaEncrypt(publicKey []byte, origData []byte) ([]byte, error) {
	var maxSize int //加密块最大长度限制
	var retErr error
	var b bytes.Buffer //直接定义一个 Buffer 变量，而不用初始化

	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)

	// /* 自己写的rsa截取拼装
	maxSize = 117
	maxi := len(origData)
	for i := 0; i < maxi; i += maxSize {
		tmpData := []byte{}
		if len(origData) > maxSize {
			tmpData = origData[:maxSize]
			origData = origData[maxSize:]
		} else {
			tmpData = origData
		}
		resByte, encryptErr := rsa.EncryptPKCS1v15(rand.Reader, pub, tmpData)
		if encryptErr != nil {
			retErr = encryptErr
		}
		b.Write(resByte)
	}
	// */

	/* 正常小于公钥长度的加密
	resByte, encryptErr := rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
	if encryptErr != nil {
		// retErr = encryptErr
	}

	return resByte, encryptErr
	*/

	return b.Bytes(), retErr

	// return base64.StdEncoding.EncodeToString(b.Bytes()), retErr

	// return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
	// return _emsa_pkcs1_v1_5_encode(origData, block.Bytes,

}

//RSA解密
func RsaDecrypt(privateKey []byte, ciphertext []byte) ([]byte, error) {
	var maxSize int //加密块最大长度限制
	var retErr error
	var b bytes.Buffer //直接定义一个 Buffer 变量，而不用初始化
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	maxSize = 128
	maxi := len(ciphertext)
	for i := 0; i < maxi; i += maxSize {
		tmpData := []byte{}
		if len(ciphertext) > maxSize {
			tmpData = ciphertext[:maxSize]
			ciphertext = ciphertext[maxSize:]
		} else {
			tmpData = ciphertext
		}
		resByte, encryptErr := rsa.DecryptPKCS1v15(rand.Reader, priv, tmpData)
		if encryptErr != nil {
			retErr = encryptErr
		}
		b.Write(resByte)
	}

	return b.Bytes(), retErr
	// return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

//RSA加签
func RsaRawSign(privateKey []byte, data []byte) ([]byte, error) {

	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error")
	}
	priInterface, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	hashed := sha1.Sum(data)

	return rsa.SignPKCS1v15(rand.Reader, priInterface, crypto.SHA1, hashed[:])
}

//RSA解签
func RsaRaw(publicKey []byte, data, sign []byte) error {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return errors.New("publicKey key error")
	}
	priInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	hashed := sha1.Sum(data)

	return rsa.VerifyPKCS1v15(priInterface.(*rsa.PublicKey), crypto.SHA1, hashed[:], sign)
}

//Rsa256RawSign RSA256加签
func Rsa256RawSign(privateKey []byte, data []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error")
	}
	priInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	if err != nil {
		fmt.Println("private key error", err)
	}
	miryanPrivateKey := priInterface.(*rsa.PrivateKey)

	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto // for simple example
	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(data)
	hashed := pssh.Sum(nil)

	// hashed := sha256.Sum256(data)

	return rsa.SignPSS(rand.Reader, miryanPrivateKey, newhash, hashed, &opts)
	// return rsa.SignPKCS1v15(rand.Reader, priInterface, crypto.SHA1, hashed[:])
}

//Rsa256Raw 验签
func Rsa256Raw(publicKey, data, sign []byte) error {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		fmt.Println("pem")
		return errors.New("publicKey key error")
	}
	priInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		fmt.Println("x509")
		return err
	}

	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto // for simple example

	newhash := crypto.SHA256
	pssh := newhash.New()
	pssh.Write(data)
	hashed := pssh.Sum(nil)

	err = rsa.VerifyPSS(priInterface.(*rsa.PublicKey), newhash, hashed, sign, &opts)

	return err
	// return nil, rsa.VerifyPKCS1v15(priInterface.(*rsa.PublicKey), crypto.SHA256, hashed[:], data)
}

//EncryptSha1 SHA1加密
func EncryptSha1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}
