package handler

import (
	"fmt"
	"net/http"
	"encoding/json"

	"webapi/gerror"
	"webapi/common/tools/tool"
	"webapi/protos/http_proto"
	"webapi/business/api_loginc"
	"webapi/common/tools/encryption"
	"bytes"
	"webapi/common/tools/httputil"
	"webapi/mem_data"
	"webapi/config"
	"github.com/mkideal/log"
)

var (
	Handlers = map[string]Command{}
	mediater = api_loginc.GetMadiator()
	cfgtools *config.Config
)

func registerCommand(cmd Command) {
	if _, found := Handlers[cmd.Name()]; found {
		panic("command " + cmd.Name() + " exsited")
	}
	Handlers[cmd.Name()] = cmd
}

type Context struct {
	Request            *http.Request
	ResponseWriter     http.ResponseWriter
	Protocol           http_proto.Request
	IsReuestDependable bool
}

type Command interface {
	Name() string
	Exec(Context) (http_proto.Response, error)
	NoLogin() bool
}

type baseCommand struct {
	name    string
	noLogin bool
}

func (cmd baseCommand) NoLogin() bool {
	return cmd.noLogin
}

func (cmd baseCommand) Name() string {
	return cmd.name
}

func JSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	json.NewEncoder(w).Encode(data)
}

func JSONResponseWithStatus(w http.ResponseWriter, data interface{}, status int) {
	w.WriteHeader(status)
	JSONResponse(w, data)
}

type HTTPResponse struct {
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data"`
}

func Init(conf *config.Config) {
	cfgtools = conf
}

func HandlerCore(w http.ResponseWriter, r *http.Request) {
	defer tool.PrintPanicStack("HandlerCore")
	// 允许跨域
	w.Header().Set("Access-Control-Allow-Origin", "*")

	for {

		if r.Method != "POST" {
			JSONResponseWithStatus(w, nil, http.StatusNotFound)
			break
		}

		cmdName := r.FormValue("cmd")

		if cmdName == "" {
			break
		}

		handler, ok := Handlers[cmdName]

		if !ok {
			var res = &http_proto.ResInfo{Status: 0, Cmd: cmdName, Description: "命令未注册"}
			JSONResponseWithStatus(w, res, http.StatusBadRequest)
			break
		}

		req, err := decryptParams(cmdName, r.FormValue("p"))

		if err != nil {
			var res = &http_proto.ResInfo{Status: 0, Cmd: cmdName, Description: err.Error()}
			JSONResponseWithStatus(w, res, http.StatusBadRequest)
			break
		}

		verifiedRequest := false
		uid := req.GetReqCommon().Uid

		if uid > 0 {
			// 用户ID总是大于0的
			userInfo, err := mem_data.LoadUserStore(uid)
			if err != nil {
				log.Info("mem_data.LoadUserStore %d sid: %v,ip:%s", uid, err, tool.GetIPAddress(r))
				var res = &http_proto.ResInfo{Status: 0, Cmd: cmdName, Description: err.Error()}
				JSONResponseWithStatus(w, res, http.StatusBadRequest)
				break
			} else {
				sid := encryption.RC4DescryptBase64(userInfo.Sid, cfgtools.PwdSalt)
				verifiedRequest = httputil.MD5Verify(sid, r.Form)
			}
		}

		if !handler.NoLogin() && !verifiedRequest {
			var res = &http_proto.ResInfo{Status: 0, Cmd: cmdName, Description: "检验错误"}
			JSONResponseWithStatus(w, res, http.StatusBadRequest)
			break
		}

		res, err := handler.Exec(Context{Request: r, ResponseWriter: w, Protocol: req, IsReuestDependable: false})

		if err != nil {
			var res = &http_proto.ResInfo{Status: 0, Cmd: cmdName, Description: err.Error()}
			JSONResponseWithStatus(w, res, http.StatusBadRequest)
			break
		}

		buf := new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(res)
		if err != nil {
			var res = &http_proto.ResInfo{Status: 0, Cmd: cmdName, Description: err.Error()}
			JSONResponseWithStatus(w, res, http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json;charset=utf-8")
			fmt.Fprint(w, encryption.RC4Base64f(buf.Bytes(), cfgtools.Rc4PwdSalt))
		}
		break
	}
}

func decryptParams(cmdName string, p string) (http_proto.Request, error) {
	req, ok := http_proto.NewRequest(cmdName)

	if !ok {
		return nil, gerror.ArgumentMissing.Append(cmdName)
	}

	if err := json.Unmarshal(tool.StrToBytes(encryption.RC4DescryptBase64(p, cfgtools.Rc4PwdSalt)), req); err != nil {
		return nil, gerror.ParseError.Append(err.Error())
	}
	return req, nil
}
