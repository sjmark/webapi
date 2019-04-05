package main

import (
	"time"
	"fmt"
	"webapi/common/tools/tool"
	"webapi/session"
	"webapi/common/netutil"
	"net"
	"webapi/config"
)

var cfg *config.Config

func startup(conf *config.Config) {
	cfg = conf
	defer tool.PrintPanicStack("startup")
	// 启动 TCP 服务
	addr := fmt.Sprintf(":%d", cfg.ImPort)

	if err := netutil.ListenAndServeTCP(addr, handleConn, true); err != nil {
		return
	}

}

// 处理新建立的连接
func handleConn(conn net.Conn) {

	id := conn.RemoteAddr().String()
	session := new(session.Session)
	reader := netutil.NewPacketReader(conn, session.OnRecvPacket)
	// 设置 3 分钟的超时
	// 如果 3 分钟没有收到客户端过来的消息
	// 连接将会断开
	reader.SetTimeout(time.Second * 180)
	session.RWSession = netutil.NewRWSession(id, cfg.ConWriteSize, reader) //s.Config.ConWriteSize
	// 启动会话
	session.Run(session.OnNewSession, session.OnQuitSession)
}
