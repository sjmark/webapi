package msg

import (
	"webapi/protos/im_proto"
)

type MsgBussiness struct{}

func (m *MsgBussiness) OnLineMsg(req *protos.OnLineMsgReq) (*protos.OnLineMsgRes, error) {
	// 消息过滤
	res := &protos.OnLineMsgRes{
		MsgType:    req.MsgType,
		MsgContent: req.MsgContent,
	}
	return res, nil
}
