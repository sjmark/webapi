package login

import (
	"webapi/protos/http_proto"
	"webapi/config"
	"webapi/common/tools/tool"
	"webapi/common/tools/encryption"
	"webapi/mem_data"
	"time"
)

type LoginLoginc struct {
	*config.Config
}

// 加密流程 服务端给客户端sid 通过sid 加密MD5Sign（param） 再通过参数拼接加密一个sign ，然后服务端校验sign 然后解密数据
func (l *LoginLoginc) Login(req *http_proto.LoginReq) (http_proto.Response, error) {
	userInfo, err := mem_data.LoadUserStoreByMobile(req.Mobile) // 后期写s

	if err != nil {
		return nil, err
	}

	if userInfo.Uid > 0 {
		if userInfo.SidExpTime <= time.Now().Unix() {
			randSid := tool.GetRandomStr(16)
			var sid = encryption.RC4Base64f(tool.StrToBytes(randSid), l.PwdSalt)
			if _, err := mem_data.LoadUserStore(userInfo.Uid, mem_data.UpdateUserSid(sid)); err != nil {
				return nil, err
			}
		}

		return http_proto.LoginRes{
			ReqCommon: http_proto.ReqCommon{Uid: userInfo.Uid},
			Sid:       encryption.RC4DescryptBase64(userInfo.Sid, l.PwdSalt),
		}, nil
	} else {

		randSid := tool.GetRandomStr(16)
		var sid = encryption.RC4Base64f(tool.StrToBytes(randSid), l.PwdSalt)
		userInfo, err := mem_data.LoadUserStore(userInfo.Uid, mem_data.NewUser(req.Mobile, sid))

		if err != nil {
			return nil, err
		}

		return http_proto.LoginRes{ReqCommon: http_proto.ReqCommon{Uid: userInfo.Uid}, Sid: randSid}, nil
	}
}
