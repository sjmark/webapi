package user

import (
	"webapi/config"
	"webapi/protos/http_proto"
	"webapi/mem_data"
)

type UserLoginc struct {
	*config.Config
}

func (u *UserLoginc) UpdateUserInfo(req *http_proto.UserUpInfoReq) (http_proto.Response, error) {
	_, err := mem_data.LoadUserStore(req.GetReqCommon().Uid, mem_data.UpdateUserInfo(req.City, req.Avatar, req.Nickname))

	if err != nil {
		return nil, err
	}

	return http_proto.CommonRes{Status: 1, Des: "success"}, nil
}
