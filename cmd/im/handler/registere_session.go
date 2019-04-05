package handler

import (
	"webapi/protos/im_proto"
	"webapi/session"
	"webapi/gerror"
)

func init() {
	session.RegisterHandler(new(protos.RegisterReq), onRegister)
}

func onRegister(ctx *session.Context) error {
	req, ok := ctx.Protocol.(*protos.RegisterReq)
	if !ok {
		return gerror.ParseError
	}
	ctx.Session.AddSession(req.UserId)
	ctx.Send(&protos.CurrencyRes{Code: 1, Desc: "ok"})
	return nil
}
