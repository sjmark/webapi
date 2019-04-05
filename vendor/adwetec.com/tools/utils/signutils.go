package utils

import (
	"io"
	"fmt"
	"sort"
	"bytes"
	"strconv"
	"net/url"
	"strings"
	"crypto/md5"
	"crypto/hmac"
	"encoding/hex"
	"encoding/json"
	"crypto/sha256"
)

// ***************************************************************************
type Values []string // 主要用于字符串数组排序

func (this Values) Len() int {
	return len(this)
}
func (this Values) Less(i, j int) bool {
	return this[i] < this[j]
}
func (this Values) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

type Filed struct {
	Key   string
	Value interface{}
}
type Fileds []*Filed

func (this Fileds) Len() int {
	return len(this)
}
func (this Fileds) Less(i, j int) bool {
	return this[i].Key < this[j].Key
}
func (this Fileds) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

// ***************************************************************************
// 函数功能: 对请求参数进行签名
// 参数说明:
// @param params:		待签名参数
// @param secret_key:	密钥
// 返 回 值: 签名字符串
// ***************************************************************************
func GenerateShoubaiHttpSign(values Values, secret_key string) string {

	params := ""

	sort.Sort(values)

	for _, val := range values {
		params += val
	}

	params += secret_key

	// fmt.Println(params)

	return Md5sum(params)
}

func Md5sum(deviceid string) string {

	sum := md5.Sum([]byte(deviceid))

	return hex.EncodeToString(sum[:])

}

func GenAuth(accesskey, secretkey, utctimestr, urlstr, method string) string {

	urlparse, err := url.Parse(urlstr)

	if err != nil {
		panic(err)
	}

	host := urlparse.Hostname()
	path := urlparse.Path

	version := "1"

	sigheaders := "host"
	expseconds := "1800"

	// ******************************************************
	// 生成签名KEY

	val := fmt.Sprintf("bce-auth-v%s/%s/%s/%s", version, accesskey, utctimestr, expseconds)

	h := hmac.New(sha256.New, []byte(secretkey))

	io.WriteString(h, val)

	signinkey := hex.EncodeToString(h.Sum(nil)) // 签名KEY

	// ******************************************************
	// URI/HEADER/REQUEST构建

	canonicaluri := path

	canonicalheaders := fmt.Sprintf("host:%s", host)

	// 待签名串
	canonicalrequest := fmt.Sprintf("%s\n%s\n\n%s", strings.ToUpper(method), canonicaluri, canonicalheaders)

	s := hmac.New(sha256.New, []byte(signinkey))

	io.WriteString(s, canonicalrequest) // 进行签名

	signature := hex.EncodeToString(s.Sum(nil)) // 签名摘要

	authorization := fmt.Sprintf("bce-auth-v%s/%s/%s/%s/%s/%s", version, accesskey, utctimestr, expseconds, sigheaders, signature)

	return authorization
}

// *******************************************************************************************
// 签名规则
// 注意: 参数值为空不参与签名并且参数名区分大小写
// Step 1: 将需要签名的数据数组中参数值为空的字段剔除得到ARRAY1
// Step 2: 将"KEY"=>"{SIGNKEY}"添加到ARRAY1中得到ARRAY2
// Step 3: 对ARRAY2进行字典排序得到ARRAY3
// Step 4: 将ARRAY3转换成JSON串得到的字符串做MD5运算得到SIGN值
// *******************************************************************************************

func GetSign(params map[string]interface{}, signkey string) string {

	// step 1.过滤参数值为空的字段

	fields := make(Fileds, 0)

	for key, value := range params {

		if value == nil {
			continue
		}

		fields = append(fields, &Filed{
			Key:   key,
			Value: value,
		})
	}

	// step 2.添加签名KEY

	fields = append(fields, &Filed{
		Key:   "key",
		Value: signkey,
	})

	// step 3.字典排序

	sort.Sort(fields)

	newparams := make(map[string]interface{})

	for _, item := range fields {
		newparams[item.Key] = item.Value
	}

	// step 4.转换为JSON串并计算MD5值
	datas, err := json.Marshal(newparams)

	if err != nil {
		return err.Error()
	}

	rs := []rune(string(datas))

	var buffer bytes.Buffer

	for _, r := range rs {

		rint := int(r)

		if rint < 128 {
			buffer.WriteString(string(r))
		} else {
			buffer.WriteString("\\u" + strconv.FormatInt(int64(rint), 16))
		}

	}

	jsonstr := strings.Replace(buffer.String(), "/", "\\/", -1) // GOLANG的JSON编码与PHP有点不一样

	return Md5sum(jsonstr)
}

//func GetSignTest(param interface{}, signkey string) string {
//
//	// step 1.过滤参数值为空的字段
//
//	newparams := make(map[string]interface{})
//
//	for key, value := range params {
//
//		if value == nil {
//			continue
//		}
//
//		newparams[key] = value
//	}
//
//	// step 2.添加签名KEY
//
//	newparams["key"] = signkey
//
//	// step 3.字典排序
//
//	// step 4.转换为JSON串并计算MD5值
//
//	datas, err := json.Marshal(newparams)
//
//	if err != nil {
//		return err.Error()
//	}
//
//	return Md5sum(string(datas))
//}
