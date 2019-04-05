package utils

import (
	"net"
	"net/http"
	_ "strings"
	_ "sync"
	"time"
)

// 默认会连接这个地址指定的TCP
// const (// 必须注释掉
//	 DefaultAddr    = "localhost:80"
//	 DefaultNetWork = "tcp"
// )

var (
	DefaultConnKeepTime = time.Duration(time.Second * 100)
)

type HttpClientPool struct { // HTTP的客户端池
	MaxPoolSize int
	DialFn      DialFunc         // 所有客户端共用
	Clients     chan *HttpClient //
}

func NewHttpClientPool(connKeepTime time.Duration, maxPoolSize int) *HttpClientPool {

	DefaultConnKeepTime = connKeepTime

	pool := &HttpClientPool{
		MaxPoolSize: maxPoolSize,
		DialFn:      DefaultDial,
		Clients:     make(chan *HttpClient, maxPoolSize),
	}

	for i := 0; i < maxPoolSize; i++ {
		pool.Clients <- NewHttpClient(pool.DialFn)
	}

	return pool
}
func (hcp *HttpClientPool) GetHttpClient() *HttpClient {
	return <-hcp.Clients
}
func (hcp *HttpClientPool) PutBack(c *HttpClient) {
	hcp.Clients <- c
}
func (hcp *HttpClientPool) SetDialFn(dialFn DialFunc) {
	hcp.DialFn = dialFn
}

type DialFunc func(network, addr string) (net.Conn, error) // 函数类型

// ***************************************************************************
type HttpClient struct { // 表示客户端
	C *http.Client // 真正的客户端
}

func NewHttpClient(dialFn DialFunc) *HttpClient {
	return &HttpClient{
		C: &http.Client{
			Transport: &http.Transport{
				Dial:              dialFn,
				DisableKeepAlives: false,
			},
		},
	}
}

// ***************************************************************************
func DefaultDial(network, addr string) (net.Conn, error) { // 默认的实现--用于创建底层的网络连接

	dial := net.Dialer{
		Timeout:   DefaultConnKeepTime,
		KeepAlive: DefaultConnKeepTime,
	}

	conn, err := dial.Dial(network, addr)

	if err != nil {
		return nil, err
	}

	return conn, err
}

// ***************************************************************************
