package http_proto

// 数据返回格式
type ResInfo struct {
	Cmd         string `json:"cmd"`
	Status      int32  `json:"status"`
	Description string `json:"description"`
}

// 所有请求协议都会带的参数
type ReqCommon struct {
	Uid int64 `json:"uid"` // 用户ID
}

// 公共返回
type CommonRes struct {
	Status int32  `json:"status"`
	Des    string `json:"des"`
}

// 登录请求
type LoginReq struct {
	Mobile string `json:"mobile"`
}

// 登录返回
type LoginRes struct {
	ReqCommon
	Sid       string `json:"sid"`
}

type UserUpInfoReq struct {
	ReqCommon
	City     string `json:"city"`
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
}
