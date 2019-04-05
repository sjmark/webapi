package im

import (
	"flag"
	//"time"

	"webapi/protos/im_proto"
	"webapi/common/packet"
	"webapi/common/netutil"

	"github.com/mkideal/log"
	pb "github.com/golang/protobuf/proto"
	"time"
	"net"
)

var flAddr = flag.String("p", "127.0.0.1:8082", "服务器地址")

func main() {
	//var origin = "http://127.0.0.1:8080/"
	conn, err := net.Dial("tcp", *flAddr)

	if err != nil {
		panic(err)
	}

	session := new(Session)
	reader := netutil.NewPacketReader(conn, session.onRecvPacket)
	session.sess = netutil.NewRWSession(conn.RemoteAddr().String(), 256, reader)
	session.sess.Run(func() {
		go func() {
			for range time.Tick(time.Second * 10) {
				heartbeat := new(protos.HeartbeatReq)
				log.Info("send heartbeat")
				send(session.sess, heartbeat)
			}
		}()

		//req := new(proto.LoginReq)
		//req.UserId=1
		//req.Username = "sujiang"
		//req.Password = "123456"
		//send(session.sess, req)
		//
		//go func() {
		//	for range time.Tick(time.Second * 3) {
		//		msgreq := new(proto.OnLineMsgReq)
		//		msgreq.ToId = []int64{2}
		//		msgreq.MsgType = 2
		//		msgreq.MsgContent = "ssssshis is test"
		//		log.Info("========send", msgreq)
		//		send(session.sess, msgreq)
		//	}
		//}()

	}, nil)
	select {}
}

func send(session netutil.Session, ptc pb.Message) {
	data, err := packet.Packet(ptc)
	if err != nil {
		log.Error("packet protocol %s error: %v", pb.MessageName(ptc), err)
		return
	}
	session.Send(netutil.BytesPacket(data))
}

type Session struct {
	sess netutil.Session
}

func (s *Session) onRecvPacket(
	data []byte,
) {
	defer func() {
		if e := recover(); e != nil {
			log.Warn("session %s onRecvPacket error: %v", s.sess.Id(), e)
		}
	}()
	typ, err := packet.DecodeType(data)
	if err != nil {
		log.Warn("解析协议类型错误: %v", err)
		s.sess.Quit()
		return
	}
	// 将协议号转换成protobuf协议号枚举并取得枚举的字符串表示
	typeName := protos.Tryout(typ).String()
	ptc := packet.New(typeName)

	if ptc == nil {
		// 不支持的协议类型
		log.Warn("不支持的协议类型: %d", typ)
		return
	}
	if err := pb.Unmarshal(data[2:], ptc); err != nil {
		log.Warn("解析消息体错误: %v", err)
		s.sess.Quit()
		return
	}
	log.Info("收到消息 %d: %v", typ, ptc)
}
