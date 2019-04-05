package main

import (
	"log"
	"fmt"
	"net/url"
	"net/http"
	"encoding/json"
	"os"
	"strconv"
	"math/rand"
	"time"
	"strings"
	"adwetec.com/spider/lib_common/redis"

	"github.com/PuerkitoBio/goquery"
)

const (
	PAGE = 40
)

var (
	xici = "http://www.xicidaili.com/wn/"
)

var cli redis.Cache

func init() {
	cli = redis.NewRedisCache("1", "127.0.0.1:6379", "", 2000)

}
func main() {
	res := GetRep("http://www.mafengwo.cn/qa/list/search?key=上海", ReturnIp())
	dom, err := goquery.NewDocumentFromResponse(res)
	fmt.Println(dom.Url, err)
	fmt.Println(dom.Html())
	return

	GetIp("local")

}

// 获取ip
func GetIp(ip string) {
	var count int

	for i := 1; i <= PAGE; i++ {

		response := GetRep(xici+strconv.Itoa(i), ip)

		if (response.StatusCode == 200) {
			dom, err := goquery.NewDocumentFromResponse(response)
			if err != nil {
				log.Fatalf("失败原因", response.StatusCode)
			}
			dom.Find("#ip_list tbody tr").Each(func(i int, context *goquery.Selection) {
				ipInfo := make(map[string][]string)
				//地址
				ip := context.Find("td").Eq(1).Text()
				//端口
				port := context.Find("td").Eq(2).Text()
				//地址
				address := context.Find("td").Eq(3).Find("a").Text()
				//匿名
				anonymous := context.Find("td").Eq(4).Text()
				//协议
				protocol := context.Find("td").Eq(5).Text()
				//存活时间
				survivalTime := context.Find("td").Eq(8).Text()
				//验证时间
				checkTime := context.Find("td").Eq(9).Text()
				ipInfo[ip] = append(ipInfo[ip], ip, port, address, anonymous, protocol, survivalTime, checkTime)
				fmt.Println(ipInfo)
				hBody, _ := json.Marshal(ipInfo[ip])

				//存入redis
				saveRedis(ip+":"+port, string(hBody))
				fmt.Println(ipInfo)
				count++
			})
		}
	}

}

func GetRep(urls string, ip string) *http.Response {

	request, _ := http.NewRequest("GET", urls, nil)
	//随机返回User-Agent 信息
	request.Header.Set("User-Agent", GetAgent())
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	request.Header.Set("Connection", "keep-alive")
	proxy, err := url.Parse(ip)
	//设置超时时间
	fmt.Printf("使用代理:%s\n", proxy)
	client := &http.Client{}

	if ip != "local" {
		client = &http.Client{
			Transport: &http.Transport{Proxy: http.ProxyURL(proxy)},
			Timeout:   time.Duration(5 * time.Second),
		}
	}

	response, err := client.Do(request)
	if err != nil || response.StatusCode != 200 {
		fmt.Println("========", err)
		GetIp(ReturnIp())
	}

	return response
}

/**
* 随机返回一个User-Agent
*/
func GetAgent() string {
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

func saveRedis(ip string, hBody string) {

	//键值对的方式存入hash
	_, err := cli.Hsetstring("ip_pool", ip, string(hBody))
	//将ip:port 存入set  方便返回随机的ip
	_, err = cli.Sadd("ippoolkey", ip)
	if err != nil {
		log.Fatalf("err:%s", err)
		os.Exit(1)
	}
}

/**
* 随机返回一个IP
*/
func ReturnIp() string {
	keys, err := cli.Smembers("ippoolkey")

	if err != nil {
		return ""
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	res, err := cli.Hget("ip_pool", keys[r.Intn(len(keys))])

	if err != nil {
		return ""
	}
	res = strings.TrimLeft(res, "[")
	res = strings.TrimRight(res, "]")

	array := strings.Split(res, ",")

	for i := 0; i < len(array); i++ {
		array[i] = strings.Trim(array[i], "\"")
	}

	host := array[0] + ":" + array[1]
	return host
}
