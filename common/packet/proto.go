package packet

import (
	"encoding/binary"
	"errors"
	"reflect"
	"strings"

	pb "github.com/golang/protobuf/proto"
)

func EncodeLength(size int) []byte {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, uint32(size))
	return data
}

func EncodeType(typ int) []byte {
	data := make([]byte, 2)
	binary.BigEndian.PutUint16(data, uint16(typ))
	return data
}

func DecodeType(data []byte) (typ int, err error) {
	if len(data) < 2 {
		err = errors.New("data length too short")
		return
	}
	typ = int(binary.BigEndian.Uint16(data))
	return
}

func Packet(ptc pb.Message) ([]byte, error) {
	typ := TypeOf(ptc)
	if typ < 0 {
		return nil, errors.New("invalid message " + pb.MessageName(ptc))
	}
	buf := append(make([]byte, 4), EncodeType(typ)...)
	buf2, err := pb.Marshal(ptc)
	if err != nil {
		return nil, err
	}
	buf = append(buf, buf2...)
	lenbuf := EncodeLength(2 + len(buf2))
	buf[0] = lenbuf[0]
	buf[1] = lenbuf[1]
	buf[2] = lenbuf[2]
	buf[3] = lenbuf[3]

	return buf, nil
}

// typeName 对应的协议不存在时, New 函数可能 panic
// 调用的地方可能需要 recover
func New(typeName string) pb.Message {
	// 根据名字获取类型
	// 加 "api." 是因为 protobuf 注册时候用名字加了包名,如果包名不是 api 请自己更换
	// 去掉 "Type" 后缀是因为枚举名后面有 Type 后缀
	// NOTE: 请保持将协议 XYZ 的协议好枚举命名为 XYZType
	// 比如 LoginReq 协议的类型枚举必须是 LoginReqType
	msgType := pb.MessageType("proto." + strings.TrimSuffix(typeName, "Type"))
	// 根据类型由反射机制创建对象
	return reflect.New(msgType.Elem()).Interface().(pb.Message)
}

// 获取协议类型
func TypeOf(ptc pb.Message) int {
	name := strings.TrimPrefix(pb.MessageName(ptc), "proto.") + "Type"
	m := pb.EnumValueMap("proto.Tryout")
	if m == nil {
		return 0
	}
	return int(m[name])
}
