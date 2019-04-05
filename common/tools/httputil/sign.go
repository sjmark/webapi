package httputil

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"net/url"
	"sort"
)

type SignFunc func([]byte) []byte

func md5sign(source []byte) []byte {
	sum := md5.Sum(source)
	return sum[:]
}

type kvpair [2]string

type pairs []kvpair

func (p pairs) Len() int           { return len(p) }
func (p pairs) Less(i, j int) bool { return p[i][0] < p[j][0] }
func (p pairs) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func join(privateKey string, args url.Values) (string, []byte) {
	var p pairs
	var sign string
	for key, values := range args {
		if key == "t" {
			continue
		}
		if len(values) == 0 {
			continue
		}
		if key == "sign" {
			sign = values[0]
			continue
		}
		p = append(p, kvpair{key, values[0]})
	}
	sort.Sort(p)
	var source bytes.Buffer
	for _, kv := range p {
		source.WriteString(kv[0])
		source.WriteString("=")
		source.WriteString(kv[1])
		source.WriteString("&")
	}
	source.WriteString(privateKey)
	return sign, source.Bytes()
}

func Sign(privateKey string, args url.Values, signFn SignFunc) string {
	_, source := join(privateKey, args)
	return base64.StdEncoding.EncodeToString(signFn(source))
}

func Verify(privateKey string, args url.Values, signFn SignFunc) bool {
	sign, source := join(privateKey, args)
	actualSign := base64.StdEncoding.EncodeToString(signFn(source))
	return sign == actualSign
}

func MD5Sign(privateKey string, args url.Values) string {
	return Sign(privateKey, args, md5sign)
}

func MD5Verify(privateKey string, args url.Values) bool {
	return Verify(privateKey, args, md5sign)
}
