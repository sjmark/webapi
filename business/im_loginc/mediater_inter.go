package im_loginc

import (
	"webapi/protos/im_proto"
)

type (
	msgInterface interface {
		OnLineMsg(req *protos.OnLineMsgReq) (*protos.OnLineMsgRes, error)
	}
)
