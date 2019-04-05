package main

import (
	"webapi/common/tools/http_tool"
	"net/url"
	"fmt"
	"webapi/protos/http_proto"
	"encoding/json"
	"webapi/common/tools/httputil"
	"webapi/common/tools/encryption"
)

func main() {
	login()
	userUp()
}

var key = "ac840da30f1a0598f1881a969426eb02"

func login() {

	// 登陆无需数据加密
	var req = &http_proto.LoginReq{
		Mobile: "xxxxxx",
	}

	s, _ := json.Marshal(req)
	var clusInfo = url.Values{}
	clusInfo.Add("cmd", "Login")
	clusInfo.Add("p", encryption.RC4Base64f(s, key))

	res, err := http_tool.HTTPDoPost(clusInfo.Encode(), "http://127.0.0.1:8287/index")

	if err != nil {
		var info = http_proto.ResInfo{}
		json.Unmarshal(res, &info)
		fmt.Println(info)
		return
	}

	// 返回数据解密
	ss := encryption.RC4DescryptBase64(string(res), key)
	var info = http_proto.LoginRes{}
	json.Unmarshal([]byte(ss), &info)

	fmt.Println(info)
}

func userUp() {

	var req = &http_proto.UserUpInfoReq{
		City:     "北京",
		Avatar:   "img.jpg",
		Nickname: "大头贵",
	}

	req.SetReqCommon(http_proto.ReqCommon{Uid: 10001})

	var sid = "aGwofpiV6RVwN3Iw"
	s, _ := json.Marshal(req)
	var clusInfo = url.Values{}
	clusInfo.Add("cmd", "UserUp")
	// 这里加密
	clusInfo.Add("p", encryption.RC4Base64f(s, key))
	sign := httputil.MD5Sign(sid, clusInfo)

	clusInfo.Add("sign", sign)

	res, err := http_tool.HTTPDoPost(clusInfo.Encode(), "http://127.0.0.1:8287/index")

	if err != nil {
		var info = http_proto.ResInfo{}
		json.Unmarshal(res, &info)
		fmt.Println(info)
		return
	}
	// 返回数据解密
	data := encryption.RC4DescryptBase64(string(res), key)

	var info = http_proto.CommonRes{}

	json.Unmarshal([]byte(data), &info)

	fmt.Println(info)
}
