# webapi

基于trie树的敏感词过滤服务， 支持grpc调用


## 安装

```
$ go get github.com/sjmark/webapi

$ cd $GOPATH/src/github.com/sjmark/webapi
$ cp webapi $GOPATH/src/
$ cd $GOPATH/src/webapi/config
$ vim conf.cnf
$ cd  $GOPATH/src/webapi/cmd/server
$ go run main.go start.go
```

```
$ cd  $GOPATH/src/webapi/cmd/filter
$ go run main.go
```

## 使用

参考 $GOPATH/src/webapi/cmd/client/server

部分代码如：

```
	// 登陆无需数据加密
	var req = &http_proto.LoginReq{
		Mobile: "158****3938",
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
	ss := encryption.RC4DescryptBase64(string(res),key)
	var info = http_proto.LoginRes{}
	json.Unmarshal([]byte(ss), &info)

	fmt.Println(info)
	
```
数据请求返回：

```
{
	"uid":10001,
	"sid":"aGwofpiV6RVwN3Iw"
}
	
```