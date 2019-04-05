package hanndler

import (
	"webapi/protos/im_proto"

	"webapi/session"
)

func init() {
	session.RegisterHandler(new(protos.HeartbeatReq), onHeartbeat)
}

func onHeartbeat(context *session.Context) error {
	// 接到心跳协议,什么都不用处理
	res := new(protos.HeartbeatRes)
	// 确认 在线、离线
	context.Send(res)
	return nil
}
