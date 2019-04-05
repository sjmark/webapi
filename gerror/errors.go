package gerror

import (
	"fmt"

	"github.com/golang/protobuf/proto"
)

//定义错误类型
var (
	Success                                  = Error(StatusSuccess, "成功")                                             // ok
	StatusFail                               = Error(StatusDefeated, "失败")                                            // fail
	NotInitialized                           = Error(StatusNotInitialized, "初始化异常")                                   // 初始化异常
	ArgumentMissing                          = Error(StatusArgumentMissing, "参数缺失")                                   // 参数缺失
	InvalidArgument                          = Error(StatusInvalidArgument, "无效参数")                                   // 无效参数
	UnexpectedError                          = Error(StatusUnexpectedError, "通用错误")                                   // 通用错误
	Unauthorized                             = Error(StatusUnauthorized, "未授权")                                       // 未授权
	DataNotAvailable                         = Error(StatusDataNotAvailable, "数据不可用")                                 // 数据不可用
	DatabaseException                        = Error(StatusDatabaseException, "数据库操作异常")                              // 数据库操作异常
	ParseError                               = Error(StatusParseError, "解析出错")                                        // 解析出错
	DataRepeatError                          = Error(StatusDataRepeatError, "重复提交数据")                                 // 该物品已激活
)

//ErrorCode 错误类型结构
type ErrorCode interface {
	error
	Code() int32
	Append(string, ...interface{}) ErrorCode
}

type errorCodeImpl struct {
	code int32
	msg  string
}

func (e errorCodeImpl) Error() string { return e.msg }

func (e errorCodeImpl) Code() int32 { return e.code }

func (e errorCodeImpl) Append(format string, args ...interface{}) ErrorCode {

	if len(args) > 0 {
		return Error(e.code, e.msg+": "+fmt.Sprintf(format, args...))
	}

	return Error(e.code, e.msg+": "+format)
}

//Error 生成错误
func Error(code int32, format string, args ...interface{}) ErrorCode {
	return &errorCodeImpl{
		code: code,
		msg:  fmt.Sprintf(format, args...),
	}
}

//Code 获取错误码
func Code(err error) int32 {

	if err == nil {
		return int32(StatusSuccess)
	}

	if e, ok := err.(ErrorCode); ok {
		return e.Code()
	}

	return int32(StatusUnexpectedError)
}

//String error to string
func String(err error) string {

	if err == nil {
		return "success"
	}

	return err.Error()
}

//Exception 异常错误
func Exception(err error) ErrorCode {

	if err == nil {
		return nil
	}

	return Error(int32(StatusUnexpectedError), err.Error())
}

func ReturnErrorRes(err error) proto.Message {
	return nil
}
