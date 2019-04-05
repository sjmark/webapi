package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"adwetec.com/tools/logrus"
	"adwetec.com/tools/protos"
	"github.com/kataras/iris/core/errors"

	"github.com/golang/protobuf/proto"
)

type HttpManager struct {
	//
	logger *logrus.Logger
	client *http.Client
}

func NewHttpManager(logger *logrus.Logger) *HttpManager {

	manager := &HttpManager{
		logger: logger,
		client: &http.Client{},
	}

	manager.client.Timeout = time.Second * 10

	return manager
}

func (this *HttpManager) SendGetHttp(urlStr string) string { // 发送GET请求

	resp, err := http.Get(urlStr)

	if err != nil {
		this.logger.Warn(fmt.Sprintf("SendGetHttp: url %s get error: %s", urlStr, err.Error()))
		return ""
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		this.logger.Warn(fmt.Sprintf("SendGetHttp: url %s readall error: %s", urlStr, err.Error()))
		return ""
	}

	return string(body)
}

func (this *HttpManager) SendPostHttpCustomHeader(urlStr, bodyParam string, header http.Header) string {
	return this.sendPostHttp(urlStr, bodyParam, header)
}
func (this *HttpManager) SendPostHttp(urlStr, bodyParam string) string {
	return this.sendPostHttp(urlStr, bodyParam, nil)
}

func (this *HttpManager) sendPostHttp(urlStr, bodyParam string, header http.Header) string { // 发送POST请求

	reader := strings.NewReader(bodyParam)

	req, err := http.NewRequest("POST", urlStr, reader)

	if err != nil {
		this.logger.Warn(fmt.Sprintf("sendPostHttp: url %s bodyparm %s new request error %s", urlStr, bodyParam, err.Error()))
		return ""
	}

	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	if header != nil {
		req.Header = header
	}

	resp, err := this.client.Do(req)

	if err != nil {
		this.logger.Warn(fmt.Sprintf("sendPostHttp: url %s bodyparm %s client do error %s", urlStr, bodyParam, err.Error()))
		return ""
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		this.logger.Warn(fmt.Sprintf("sendPostHttp: url %s bodyparm %s read all error %s", urlStr, bodyParam, err.Error()))
		return ""
	}

	return string(body)
}

//*********************************************************************************
func SendGetHttp(urlStr string) (string, error) { // 发送GET请求

	resp, err := http.Get(urlStr)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil
}
func SendPostHttp(url string) string {

	request, err := http.NewRequest("POST", url, nil)

	if err != nil {
		return err.Error()
	}

	//	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//	request.Header.Set("Connection", "Keep-Alive")

	resp, err := http.DefaultClient.Do(request)

	if err != nil { // 出错时需要重试
		return err.Error()
	}

	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err.Error()
	}

	// {"body":{"header":{"code":20000,"msg":"请求成功！"},"entity":null}}

	//fmt.Println(string(body))
	//
	//var data map[string]interface{}
	//
	//json.Unmarshal(body, &data)
	//
	//fmt.Println(data["errno"])
	//fmt.Println(data["data"])

	return string(body)

}
func SendPostHttpAudit(url string, datas string) (string, error) {

	body := ioutil.NopCloser(strings.NewReader(datas))

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", err
	}

	req.Header.Set("content-type", "application/json")

	resp, err := client.Do(req) //发送

	if err != nil {
		return "", err
	}

	defer resp.Body.Close() // 一定要关闭resp.Body

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(data), nil

}
func SendPostHttpFeed(url, utc_time_str, auth, datas string) (string, error) {

	body := ioutil.NopCloser(strings.NewReader(datas))

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", err
	}

	req.Header.Set("accept-encoding", "gzip, deflate")
	req.Header.Set("host", "sem.baidubce.com")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-bce-date", utc_time_str)
	req.Header.Set("authorization", auth)
	req.Header.Set("accept", "*/*")

	// fmt.Printf("%+v\n", req) // 看下发送的结构

	resp, err := client.Do(req) //发送
	if err != nil {
		return "", err
	}

	defer resp.Body.Close() //一定要关闭resp.Body

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil

}
func SendPostHttpSem(url, utc_time_str, auth, datas string) (string, error) {

	body := ioutil.NopCloser(strings.NewReader(datas))

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", err
	}

	req.Header.Set("accept-encoding", "gzip, deflate")
	req.Header.Set("host", "sem.baidubce.com")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-bce-date", utc_time_str)
	req.Header.Set("authorization", auth)
	req.Header.Set("accept", "*/*")

	// fmt.Printf("%+v\n", req) // 看下发送的结构

	resp, err := client.Do(req) //发送
	if err != nil {
		return "", err
	}

	defer resp.Body.Close() //一定要关闭resp.Body

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil

}
func SendPostHttpHuichuan(url, datas string) (string, error) {

	body := ioutil.NopCloser(strings.NewReader(datas))

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req) //发送

	if err != nil {
		return "", err
	}

	defer resp.Body.Close() //一定要关闭resp.Body

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil

}
func SendPostHttpHuichuan2(url, datas string) (string, error) {

	body := ioutil.NopCloser(strings.NewReader(datas))

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "multipart/form-data")

	resp, err := client.Do(req) //发送

	if err != nil {
		return "", err
	}

	defer resp.Body.Close() //一定要关闭resp.Body

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil

}
func SendPostHttpShenma(url, datas string) (string, error) {

	body := ioutil.NopCloser(strings.NewReader(datas))

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req) //发送

	if err != nil {
		return "", err
	}

	defer resp.Body.Close() //一定要关闭resp.Body

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil

}
func HttpRequest(jsonBody interface{}, path string) {

	bytesData, err := json.Marshal(jsonBody)

	if err != nil {
		fmt.Println("json.mashal error", err.Error())
		return
	}

	reader := bytes.NewReader(bytesData)

	url := "http://localhost:3111/api/wetec-service/" + path

	request, err := http.NewRequest("POST", url, reader)

	if err != nil {
		fmt.Println("http.NewRequest error", err.Error())
		return
	}

	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MjQ1NTY1MzMsImlhdCI6MTUyNDU1NjIzMywicm9sZWlkIjoxLCJ1c2VyaWQiOjF9.iZjSP5OE9m_E_WfdbBmcKgLZ1ufSGYC-P3vnUZ8yp7I")

	client := http.Client{}

	resp, err := client.Do(request)

	if err != nil {
		fmt.Println("client.Do error", err.Error())
		return
	}

	respBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("ioutil.ReadAll error", err.Error())
		return
	}

	var out bytes.Buffer

	err = json.Indent(&out, respBytes, "", "\t")

	if err != nil {
		fmt.Println("json.Indent error", err.Error())
		return
	}

	out.WriteTo(os.Stdout)
}

func HttpRequestProto(protostr proto.Message, path string) (*protos.PdbBidResponse, error) {

	body, err := proto.Marshal(protostr)

	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(body)

	// url := "http://bidupdbreq.adwintech.com/" + path
	url := "http://localhost:41111/" + path

	request, err := http.NewRequest("POST", url, reader)

	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/octet-stream")
	request.Header.Set("Connection", "keep-alive")

	client := http.Client{}

	resp, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	pbidres := &protos.PdbBidResponse{}

	err = proto.Unmarshal(respBytes, pbidres)

	if err != nil {
		return nil, err
	}

	return pbidres, nil
}
func HttpPostJson(url, datas string) (string, error) {

	body := ioutil.NopCloser(strings.NewReader(datas))

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", err
	}

	req.Header.Set("content-type", "application/json")

	resp, err := client.Do(req) //发送
	if err != nil {
		return "", err
	}

	defer resp.Body.Close() //一定要关闭resp.Body

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
func httpPost() {

	// 要求第二个参数必须下面这种否则POST参数无法传递

	resp, err := http.Post("http://www.01happy.com/demo/accept.php", "application/x-www-form-urlencoded", strings.NewReader("name=cjb"))

	if err != nil {

		fmt.Println(err)

	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}
func httpPostForm() {

	resp, err := http.PostForm("http://www.01happy.com/demo/accept.php", url.Values{"key": {"Value"}, "id": {"123"}})

	if err != nil {
		// handle error
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		// handle error
	}

	fmt.Println(string(body))

}
func httpDo() { // 可以设置头参数\COOKIE之类的数据

	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://www.01happy.com/demo/accept.php", strings.NewReader("name=cjb"))

	if err != nil {
		// handle error
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "name=anny")

	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}

//腾讯
func SendPostHttpTx(url string, datas string) (string, error) {

	body := ioutil.NopCloser(strings.NewReader(datas))

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", err
	}

	req.Header.Set("content-type", "application/json")

	resp, err := client.Do(req) //发送
	if err != nil {
		return "", err
	}

	defer resp.Body.Close() //一定要关闭resp.Body

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil

}

//腾讯postfrom形式文件提交
func HttpPostFormTx(url string, otherParams map[string]string, fileName string) (string, error) {

	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)

	// file
	fileWriter, err := bodyWriter.CreateFormFile("file", fileName)
	if err != nil {
		return "", err
	}
	file, err := os.Open(fileName)
	if err != nil {
		return "", nil
	}

	defer file.Close()
	io.Copy(fileWriter, file)

	//同文件一起上传的其他参数
	extraParams := otherParams
	for key, value := range extraParams {
		err := bodyWriter.WriteField(key, value)
		if err != nil {
			return "", err
		}
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(url, contentType, bodyBuffer)
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
func SendGetHttpToutiao(urlStr string, token string) (string, error) { // 发送GET请求附带Token
	client := &http.Client{}

	req, err := http.NewRequest("GET", urlStr, nil)
	req.Header.Set("Access-Token", token)
	resp, err := client.Do(req) //发送
	if err != nil {
		return "", err
	}

	defer resp.Body.Close() //一定要关闭resp.Body

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func SendPostHttpTt(url, datas string, token string) (string, error) {

	body := ioutil.NopCloser(strings.NewReader(datas))

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, body)

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Access-Token", token)

	resp, err := client.Do(req) //发送

	if err != nil {
		return "", err
	}

	defer resp.Body.Close() //一定要关闭resp.Body

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil

}

func SendPostHttpTtNoToken(url, datas string) (string, error) {

	body := ioutil.NopCloser(strings.NewReader(datas))

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req) //发送

	if err != nil {
		return "", err
	}

	defer resp.Body.Close() //一定要关闭resp.Body

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil

}

// 发送GET请求并将返回值保存在指定路径下
func GetHttpToPath(path string, url string) error {

	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	file, err := os.Create(path)

	if err != nil {
		return err
	}

	io.Copy(file, resp.Body)

	return nil
}

func SendPostHttp2(url string, datas string) (string, error) {

	body := ioutil.NopCloser(strings.NewReader(datas))

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req) //发送

	if err != nil {
		return "", err
	}

	defer resp.Body.Close() //一定要关闭resp.Body

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil

}


//接口调用公共方法，请求post
func SendHttp(params string, method string, sendurl string) ([]byte, error) {

	if method == "POST" {
		client := &http.Client{}
		req, err := http.NewRequest(method, sendurl, bytes.NewBuffer([]byte(params)))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp, err := client.Do(req)

		if err != nil {
			return []byte{}, err
		}

		if resp.StatusCode != 200 {

			return []byte{}, errors.New(fmt.Sprintf("resp.StatusCode:%v", resp.StatusCode))
		} else {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return []byte{}, err
			}
			return body, nil
		}
	}
	return []byte{}, errors.New("Illegal operation")
}