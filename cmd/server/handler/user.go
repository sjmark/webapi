package handler

import (
	"webapi/protos/http_proto"
	"webapi/gerror"
	"webapi/common/tools/tool"
	"context"
	"webapi/protos/filter_proto"
	"webapi/common/filter"
)

type userCommand struct {
	baseCommand
}

func init() {
	registerCommand(newUserCommand())
}

func newUserCommand() *userCommand {
	cmd := &userCommand{}
	cmd.name = http_proto.UserUpType
	return cmd
}

func (cmd *userCommand) Exec(ctx Context) (res http_proto.Response, err error) {
	ptc := ctx.Protocol.(*http_proto.UserUpInfoReq)

	if filter.CityIsExists(ptc.City) {
		err = gerror.UnexpectedError.Append("城市不存在")
		return
	}

	if str := tool.ValidateStr(ptc.City); len(str) == 0 {
		err = gerror.UnexpectedError.Append("请输入合理的城市")
		return
	}

	if str := tool.ValidateStr(ptc.Nickname); len(str) == 0 {
		err = gerror.UnexpectedError.Append("请输入合理的昵称")
		return
	}

	t := new(protos.Text)
	t.Text = ptc.Nickname

	keys, err1 := protos.NewWordFilterClient(cfgtools.FilterConn).FindKeyWords(context.Background(), t)

	if err1 != nil {
		err = err1
		return
	}

	if len(keys.KeyWords) > 0 {
		err = gerror.UnexpectedError.Append("名称存在过敏词汇")
		return
	}

	if !tool.ValidateImgStr(ptc.Avatar) {
		err = gerror.UnexpectedError.Append("请输入合理的图片")
		return
	}

	return mediater.UserMod.UpdateUserInfo(ptc)
}
