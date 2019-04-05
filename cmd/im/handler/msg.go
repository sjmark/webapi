package handler

import (
	"webapi/session"
	"webapi/protos/im_proto"
	"webapi/business/im_loginc"
)

func init() {
	session.RegisterHandler(new(protos.OnLineMsgReq), onLineMsg)
}

func onLineMsg(ctx *session.Context) error {
	req := ctx.Protocol.(*protos.OnLineMsgReq)
	res, err := im_loginc.GetMadiator().MsgMod.OnLineMsg(req)

	if err != nil {
		return err
	}
	ctx.Broadcast(req.ToId, res)
	return nil
}
