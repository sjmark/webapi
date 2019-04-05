package typeconv

import (
	"fmt"
	"strconv"
)

//Unused 无效
const Unused = 0

//SetInt64FromStr int64 -> string
func SetInt64FromStr(ptr *int64, s string) error {
	i, err := strconv.ParseInt(s, 0, 64)
	if err == nil {
		*ptr = i
	}
	return err
}

//SetInt32FromStr int32 -> string
func SetInt32FromStr(ptr *int32, s string) error {
	i, err := strconv.ParseInt(s, 0, 64)
	if err == nil {
		*ptr = int32(i)
	}
	return err
}

//SetInt16FromStr int16 -> string
func SetInt16FromStr(ptr *int16, s string) error {
	i, err := strconv.ParseInt(s, 0, 64)
	if err == nil {
		*ptr = int16(i)
	}
	return err
}

//SetInt8FromStr int8 -> string
func SetInt8FromStr(ptr *int8, s string) error {
	i, err := strconv.ParseInt(s, 0, 64)
	if err == nil {
		*ptr = int8(i)
	}
	return err
}

//SetIntFromStr int -> string
func SetIntFromStr(ptr *int, s string) error {
	i, err := strconv.ParseInt(s, 0, 64)
	if err == nil {
		*ptr = int(i)
	}
	return err
}

//SetUint64FromStr uint64 -> string
func SetUint64FromStr(ptr *uint64, s string) error {
	i, err := strconv.ParseUint(s, 0, 64)
	if err == nil {
		*ptr = i
	}
	return err
}

//SetUint32FromStr uint32 -> string
func SetUint32FromStr(ptr *uint32, s string) error {
	i, err := strconv.ParseUint(s, 0, 64)
	if err == nil {
		*ptr = uint32(i)
	}
	return err
}

//SetUint16FromStr Uint16 -> string
func SetUint16FromStr(ptr *uint16, s string) error {
	i, err := strconv.ParseUint(s, 0, 64)
	if err == nil {
		*ptr = uint16(i)
	}
	return err
}

//SetUint8FromStr uint8 -> string
func SetUint8FromStr(ptr *uint8, s string) error {
	i, err := strconv.ParseUint(s, 0, 64)
	if err == nil {
		*ptr = uint8(i)
	}
	return err
}

//SetUintFromStr uint -> string
func SetUintFromStr(ptr *uint, s string) error {
	i, err := strconv.ParseUint(s, 0, 64)
	if err == nil {
		*ptr = uint(i)
	}
	return err
}

//SetFloat32FromStr float32 -> string
func SetFloat32FromStr(ptr *float32, s string) error {
	f, err := strconv.ParseFloat(s, 64)
	if err == nil {
		*ptr = float32(f)
	}
	return err
}

//SetFloat64FromStr float64 -> string
func SetFloat64FromStr(ptr *float64, s string) error {
	f, err := strconv.ParseFloat(s, 64)
	if err == nil {
		*ptr = float64(f)
	}
	return err
}

//SetBoolFromStr bool -> string
func SetBoolFromStr(ptr *bool, s string) error {
	if s == "" {
		*ptr = false
		return nil
	}
	b, err := strconv.ParseBool(s)
	if err == nil {
		*ptr = b
	}
	return err
}
func ObjToString(obj interface{}) string {
	return fmt.Sprint(obj)
}

// ParseFixedWidthInt64 解析固定长度整数
// 1) 比如 "123", "00123" 均解析为 123
func ParseFixedWidthInt64(str string) (int64, error) {
	tmp := int64(0)
	for i := 0; i < len(str); i++ {
		b := str[i] - '0'
		if b < 0 || b > 9 {
			return 0, fmt.Errorf("%d th byte(%c) not a digit", i, str[i])
		}
		tmp = tmp*10 + int64(b)
	}
	return tmp, nil
}

// ParseFixedWidthInt 解析固定长度整数
func ParseFixedWidthInt(str string) (int, error) {
	i, err := ParseFixedWidthInt64(str)
	return int(i), err
}
