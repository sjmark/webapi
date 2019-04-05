package handler

import (
	"webapi/protos/http_proto"
	"webapi/common/tools/tool"
	"webapi/gerror"
)

type loginCommand struct {
	baseCommand
}

func init() {
	registerCommand(newLoginCommand())
}

func newLoginCommand() *loginCommand {
	cmd := &loginCommand{}
	cmd.name = http_proto.LoginType
	cmd.noLogin = true
	return cmd
}

func (cmd *loginCommand) Exec(ctx Context) (res http_proto.Response, err error) {
	ptc := ctx.Protocol.(*http_proto.LoginReq)

	if ok := tool.MobileNumCheck(ptc.Mobile); !ok {
		err = gerror.UnexpectedError.Append("请输入合理的手机号")
		return
	}
	return mediater.LoginMod.Login(ptc)
}
