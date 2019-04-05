package session

import (
	"sync/atomic"

	"webapi/common/packet"
	"webapi/protos/im_proto"
	"webapi/common/netutil"

	"github.com/mkideal/log"
	pb "github.com/golang/protobuf/proto"
)

type Context struct {
	Session  *Session
	Protocol pb.Message
}

// 发协议给当前 context 下的 session 对应的客户端
func (context *Context) Send(ptc pb.Message) {
	data, err := packet.Packet(ptc)
	if err != nil {
		log.Error("packet protocol error: %v", err)
		return
	}
	context.Session.Send(netutil.BytesPacket(data))
}

// 广播消息给 recerivers 所代表的用户ID
func (context *Context) Broadcast(receivers []int64, ptc pb.Message) {
	if len(receivers) == 0 {
		return
	}

	data, err := packet.Packet(ptc)
	if err != nil {
		log.Error("packet protocol error: %v", err)
		return
	}
	context.Session.Broadcast(receivers, netutil.BytesPacket(data))
}

// 广播消息给所有用户
func (context *Context) BroadcastAll(ptc pb.Message) {
	data, err := packet.Packet(ptc)
	if err != nil {
		log.Error("packet protocol error: %v", err)
		return
	}
	context.Session.Broadcast(nil, netutil.BytesPacket(data))
}

type Session struct {
	//server *Server
	*netutil.RWSession
	userId int64
}

func (sess *Session) OnNewSession() {
	// 新连接建立后自动回调这个函数
	// 如果没有需要什么要做的,这个函数留空就行
}

func (sess *Session) OnQuitSession() {
	// 连接断开后会最终自动调用这个函数
	// 可以在这里处理玩家连接断开后的逻辑
	// 清空session
	LoadServerStore(RemoveSession(sess.GetUserId()))
}

func (sess *Session) OnRecvPacket(data []byte) {
	defer func() {
		if e := recover(); e != nil {
			log.Warn("session %s onRecvPacket error: %v", sess.Id(), e)
		}
	}()

	typ, err := packet.DecodeType(data)

	if err != nil {
		// 解析数据失败,打印日志然后关掉连接
		log.Warn("解析协议类型错误: %v", err)
		sess.Quit()
		return
	}

	// 将协议号转换成protobuf协议号枚举并取得枚举的字符串表示
	typeName := protos.Tryout(typ).String()
	ptc := packet.New(typeName)
	if ptc == nil {
		// 不支持的协议类型,打印日志
		log.Warn("不支持的协议类型: %d", typ)
		return
	}
	if err := pb.Unmarshal(data[2:], ptc); err != nil {
		log.Warn("解析消息体错误: %v", err)
		sess.Quit()
		return
	}
	if sess.userId == 0 {
		// userId 为 0 表示是为登陆状态
		// 未登陆时只允许一部分协议
		mtype := protos.Tryout(typ)
		switch mtype {
		case protos.Tryout_HeartbeatReqType:
		case protos.Tryout_RegisterReqType:
		default:
			// 其他协议必须要求登陆后才可以发
			log.Warn("登陆后才可以发消息 %d", typ)
			return
		}
	}
	// 派发消息进行处理
	dispatch(sess, ptc)
}

func dispatch(session *Session, ptc pb.Message) bool {

	h, ok := Handlers[packet.TypeOf(ptc)]
	if ok {
		h(&Context{
			Session:  session,
			Protocol: ptc,
		})
	}

	return ok
}

func (sess *Session) GetUserId() int64 {
	return atomic.LoadInt64(&sess.userId)
}

func (sess *Session) AddSession(userId int64) {
	atomic.StoreInt64(&sess.userId, userId)
	LoadServerStore(AddSession(userId, sess))
}

func (sess *Session) SetUserId(userId int64) {
	atomic.StoreInt64(&sess.userId, userId)
}

func (sess *Session) Broadcast(receivers []int64, data netutil.Packet) {

	server := LoadServerStore()

	server.locker.Lock()
	defer server.locker.Unlock()

	if len(receivers) > 0 {
		for _, receiver := range receivers {
			if session, ok := server.sessions[receiver]; ok {
				session.Send(data)
			}
		}
	} else {
		for _, session := range server.sessions {
			session.Send(data)
		}
	}
}

var Handlers = make(map[int]Handler)

type Handler func(*Context) error

func RegisterHandler(ptc pb.Message, h Handler) {
	Handlers[packet.TypeOf(ptc)] = h
}
