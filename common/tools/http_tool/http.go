package http_tool

import (
	"fmt"
	"time"
	"bytes"
	"net/url"
	"net/http"
	"io/ioutil"
	"math/rand"
	"encoding/json"
	"net"
	"strings"
	"context"
)

//HTTPDoPost 接口调用公共方法，请求post
func HTTPDoPost(params string, sendurl string) ([]byte, error) {
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, time.Second)
	req, err := http.NewRequest(http.MethodPost, sendurl, strings.NewReader(params))

	if err != nil {
		return []byte{}, err
	}

	req = req.WithContext(ctx)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return []byte{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return body, fmt.Errorf("resp.StatusCode:%v", resp.StatusCode)
	}

	return body, nil
}

// post 的另外一种加超时方式
func HttpPost(addr string, data url.Values, timeout time.Duration) ([]byte, error) {

	resp, err := getClient(timeout).PostForm(addr, data)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("resp.StatusCode:%v", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

// 加入超时
func GetHttp(sendurl string, timeout time.Duration) ([]byte, error) {

	req, err := getClient(timeout).Get(sendurl)
	defer req.Body.Close()

	if err != nil {
		return []byte{}, err
	}

	if req.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("resp.StatusCode:%v", req.StatusCode)
	}

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

func getClient(timeout time.Duration) *http.Client {
	if timeout < 0 {
		timeout = 3
	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netW, addr string) (net.Conn, error) {
				//设置建立连接超时
				conn, err := net.DialTimeout(netW, addr, time.Second*timeout)

				if err != nil {
					return nil, err
				}

				//设置发送接受数据超时
				conn.SetDeadline(time.Now().Add(time.Second * timeout))
				return conn, nil
			},
		},
	}
	return client
}

func SendNormalHttp(url2 string, method string, params map[string]interface{}) (*http.Response, error) {
	bytesData, err := json.Marshal(params)

	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, url2, bytes.NewReader(bytesData))
	request.Header.Set("Transfer-Encoding", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	request.Header.Set("Cache-Control", "")
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Pragma", "no-cache")
	request.Header.Set("Vary", "Accept-Encoding")
	request.Header.Set("User-Agent", "Baiduspider")

	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	resp, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("404 not found url:%v", url2)
	}

	return resp, nil
}

func sendHttp(url2 string, method string) (res []byte, urlPath string, err error) {
	proxy, err := url.Parse("http://119.176.66.90:9999")

	client := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxy)},
	}

	request, err := http.NewRequest(method, url2, nil)
	request.Header.Set("Transfer-Encoding", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	request.Header.Set("Cache-Control", "")
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Pragma", "no-cache")
	request.Header.Set("Vary", "Accept-Encoding")
	request.Header.Set("User-Agent", "Baiduspider")

	if err != nil {
		return
	}

	resp, err := client.Do(request)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("404 not found url:%v", url2)
		return
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	urlPath = resp.Request.URL.Path
	res = respBytes
	return
}

func SendHttpToBytes(url2 string, method string, params map[string]interface{}) (res []byte, urlPath string, err error) {
	proxy, err := url.Parse("http://119.176.66.90:9999")

	client := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxy)},
	}

	request, err := http.NewRequest(method, url2, nil)
	request.Header.Set("Transfer-Encoding", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	request.Header.Set("Cache-Control", "")
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Pragma", "no-cache")
	request.Header.Set("Vary", "Accept-Encoding")
	request.Header.Set("User-Agent", "Baiduspider")

	if err != nil {
		return
	}

	resp, err := client.Do(request)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("404 not found url:%v", url2)
		return
	}

	respBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return
	}

	urlPath = resp.Request.URL.Path
	res = respBytes
	return
}

func SendMFWHttp(sendUrl, ip string, method string, params map[string]interface{}) (*http.Response, error) {
	bytesData, err := json.Marshal(params)
	reader := bytes.NewReader(bytesData)

	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, sendUrl, reader)

	if err != nil {
		return nil, err
	}

	request.Header.Set("Pragma", "no-cache")
	request.Header.Set("Host", "www.mafengwo.cn")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Cache-Control", "max-age=0")
	request.Header.Set("Upgrade-Insecure-Requests", "1")
	//request.Header.Set("Proxy-Connection:", "keep-alive")
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Transfer-Encoding", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0,gzip, deflate")
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.86 Safari/537.36")
	//request.Header.Set("User-Agent", GetRandAgent())

	proxy, _ := url.Parse(ip)

	var client = &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxy)},
		Timeout:   time.Duration(5 * time.Second),
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("not found 404")
	}

	return resp, nil
}

/**
* 随机返回一个User-Agent
*/
func GetRandAgent() string {
	agent := [...]string{
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:50.0) Gecko/20100101 Firefox/50.0",
		"Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; en) Presto/2.8.131 Version/11.11",
		"Opera/9.80 (Windows NT 6.1; U; en) Presto/2.8.131 Version/11.11",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; 360SE)",
		"Mozilla/5.0 (Windows NT 6.1; rv:2.0.1) Gecko/20100101 Firefox/4.0.1",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; The World)",
		"User-Agent,Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_8; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"User-Agent, Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Maxthon 2.0)",
		"User-Agent,Mozilla/5.0 (Windows; U; Windows NT 6.1; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
	}
	return agent[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(agent))]
}
