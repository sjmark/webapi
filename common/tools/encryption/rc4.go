package encryption

import (
	"crypto/rc4"
	"encoding/base64"
	"webapi/common/tools/tool"
)

func RC4Base64f(srctTmp []byte, key string) string {
	cl, _ := rc4.NewCipher(tool.StrToBytes(key))
	dst := make([]byte, len(srctTmp))
	cl.XORKeyStream(dst, srctTmp)
	return base64.StdEncoding.EncodeToString(dst)
}

func RC4DescryptBase64(p, keystr string) string {
	strByte, err := base64.StdEncoding.DecodeString(p)
	if err != nil {
		return ""
	}
	ct, err := rc4.NewCipher(tool.StrToBytes(keystr))

	if err != nil {
		return ""
	}

	dst := make([]byte, len(strByte))

	ct.XORKeyStream(dst, strByte)
	return tool.BytesToStr(dst)
}
