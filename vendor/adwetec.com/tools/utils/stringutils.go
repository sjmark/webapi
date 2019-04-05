package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unsafe"

	"github.com/axgle/mahonia"
	"github.com/mozillazg/go-pinyin"
	"strconv"
)

func GetUnicodeString(raw string) string {

	textQuoted := strconv.QuoteToASCII(raw)

	return textQuoted[1: len(textQuoted)-1]
}

func String(obj interface{}) string {

	b, err := json.Marshal(obj)

	if err != nil {
		return fmt.Sprintf("%+v", obj)
	}

	var out bytes.Buffer

	err = json.Indent(&out, b, "", "    ")

	if err != nil {
		return fmt.Sprintf("%+v", obj)
	}

	return out.String()

}
func JsonString(obj interface{}) string {

	b, err := json.Marshal(obj)

	if err != nil {
		return fmt.Sprintf("%+v", obj)
	}

	return string(b)

}

// string 转为byte 高效
func StrToBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// byte 转为string 高效
func BytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// 数组截取
func SplitArr(
	array []int64,
	pageIndex, offset, arrayLen int,
) []int64 {

	if offset > arrayLen {
		return array[pageIndex:]
	}

	return array[pageIndex:offset]
}

// 汉字转拼音
var args *pinyin.Args

func newPinyinArgs() pinyin.Args {
	if args == nil {
		a := pinyin.NewArgs()
		args = &a
	}
	return *args
}

/**
 * @note 将汉字转成拼音
 * @params s 汉字字符串 *string
 */
func HanziToPinyin(str string) string {
	return strings.Join(pinyin.LazyPinyin(str, newPinyinArgs()), "")
}

/**
 * @note 重新编码防止乱码！split方法可能导致乱码
 *
 */
func ConvertToString(src string, srcCode string, tagCode string) string { // 公共方法

	srcCoder := mahonia.NewDecoder(srcCode)
	tagCoder := mahonia.NewDecoder(tagCode)

	srcResult := srcCoder.ConvertString(src)

	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)

	result := string(cdata)

	return result
}

/**
 * @note 获取随机字符串,l代表长度要10的就输入10.8就输入8
 *
 */
func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func HanzhiToUnicode(stext string) string {

	quoted := strconv.QuoteToASCII(stext)

	unquoted := quoted[1: len(quoted)-1]

	return unquoted
}
func UnicodeToHanzhi(unicode string) (string, error) {

	sunicodev := strings.Split(unicode, "\\u")

	var context string

	for _, v := range sunicodev {

		if len(v) < 1 {
			continue
		}

		temp, err := strconv.ParseInt(v, 16, 32)

		if err != nil {
			return "", err
		}

		context += fmt.Sprintf("%c", temp)
	}

	return context, nil
}
