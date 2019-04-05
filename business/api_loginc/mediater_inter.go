package api_loginc

import (
	"webapi/protos/http_proto"
)

type (
	loginInterface interface {
		Login(req *http_proto.LoginReq) (http_proto.Response,error)
	}

	userInterface interface {
		UpdateUserInfo(req *http_proto.UserUpInfoReq) (http_proto.Response, error)
	}
)
