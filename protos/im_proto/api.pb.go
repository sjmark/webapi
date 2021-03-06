// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api.proto

package protos

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Tryout int32

const (
	Tryout_EchoType         Tryout = 0
	Tryout_HeartbeatReqType Tryout = 101
	Tryout_HeartbeatResType Tryout = 102
	Tryout_RegisterReqType  Tryout = 103
	Tryout_OnLineMsgReqType Tryout = 1001
	Tryout_OnLineMsgResType Tryout = 1002
	Tryout_CurrencyResType  Tryout = 3000
)

var Tryout_name = map[int32]string{
	0:    "EchoType",
	101:  "HeartbeatReqType",
	102:  "HeartbeatResType",
	103:  "RegisterReqType",
	1001: "OnLineMsgReqType",
	1002: "OnLineMsgResType",
	3000: "CurrencyResType",
}
var Tryout_value = map[string]int32{
	"EchoType":         0,
	"HeartbeatReqType": 101,
	"HeartbeatResType": 102,
	"RegisterReqType":  103,
	"OnLineMsgReqType": 1001,
	"OnLineMsgResType": 1002,
	"CurrencyResType":  3000,
}

func (x Tryout) String() string {
	return proto.EnumName(Tryout_name, int32(x))
}
func (Tryout) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_api_698b2458fee38050, []int{0}
}

type HeartbeatReq struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HeartbeatReq) Reset()         { *m = HeartbeatReq{} }
func (m *HeartbeatReq) String() string { return proto.CompactTextString(m) }
func (*HeartbeatReq) ProtoMessage()    {}
func (*HeartbeatReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_api_698b2458fee38050, []int{0}
}
func (m *HeartbeatReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HeartbeatReq.Unmarshal(m, b)
}
func (m *HeartbeatReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HeartbeatReq.Marshal(b, m, deterministic)
}
func (dst *HeartbeatReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HeartbeatReq.Merge(dst, src)
}
func (m *HeartbeatReq) XXX_Size() int {
	return xxx_messageInfo_HeartbeatReq.Size(m)
}
func (m *HeartbeatReq) XXX_DiscardUnknown() {
	xxx_messageInfo_HeartbeatReq.DiscardUnknown(m)
}

var xxx_messageInfo_HeartbeatReq proto.InternalMessageInfo

type HeartbeatRes struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HeartbeatRes) Reset()         { *m = HeartbeatRes{} }
func (m *HeartbeatRes) String() string { return proto.CompactTextString(m) }
func (*HeartbeatRes) ProtoMessage()    {}
func (*HeartbeatRes) Descriptor() ([]byte, []int) {
	return fileDescriptor_api_698b2458fee38050, []int{1}
}
func (m *HeartbeatRes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HeartbeatRes.Unmarshal(m, b)
}
func (m *HeartbeatRes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HeartbeatRes.Marshal(b, m, deterministic)
}
func (dst *HeartbeatRes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HeartbeatRes.Merge(dst, src)
}
func (m *HeartbeatRes) XXX_Size() int {
	return xxx_messageInfo_HeartbeatRes.Size(m)
}
func (m *HeartbeatRes) XXX_DiscardUnknown() {
	xxx_messageInfo_HeartbeatRes.DiscardUnknown(m)
}

var xxx_messageInfo_HeartbeatRes proto.InternalMessageInfo

type RegisterReq struct {
	UserId               int64    `protobuf:"varint,1,opt,name=user_id,json=userId" json:"user_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterReq) Reset()         { *m = RegisterReq{} }
func (m *RegisterReq) String() string { return proto.CompactTextString(m) }
func (*RegisterReq) ProtoMessage()    {}
func (*RegisterReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_api_698b2458fee38050, []int{2}
}
func (m *RegisterReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterReq.Unmarshal(m, b)
}
func (m *RegisterReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterReq.Marshal(b, m, deterministic)
}
func (dst *RegisterReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterReq.Merge(dst, src)
}
func (m *RegisterReq) XXX_Size() int {
	return xxx_messageInfo_RegisterReq.Size(m)
}
func (m *RegisterReq) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterReq.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterReq proto.InternalMessageInfo

func (m *RegisterReq) GetUserId() int64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

// 通用返回
type CurrencyRes struct {
	Code                 uint32   `protobuf:"varint,1,opt,name=code" json:"code,omitempty"`
	Desc                 string   `protobuf:"bytes,2,opt,name=desc" json:"desc,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CurrencyRes) Reset()         { *m = CurrencyRes{} }
func (m *CurrencyRes) String() string { return proto.CompactTextString(m) }
func (*CurrencyRes) ProtoMessage()    {}
func (*CurrencyRes) Descriptor() ([]byte, []int) {
	return fileDescriptor_api_698b2458fee38050, []int{3}
}
func (m *CurrencyRes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CurrencyRes.Unmarshal(m, b)
}
func (m *CurrencyRes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CurrencyRes.Marshal(b, m, deterministic)
}
func (dst *CurrencyRes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CurrencyRes.Merge(dst, src)
}
func (m *CurrencyRes) XXX_Size() int {
	return xxx_messageInfo_CurrencyRes.Size(m)
}
func (m *CurrencyRes) XXX_DiscardUnknown() {
	xxx_messageInfo_CurrencyRes.DiscardUnknown(m)
}

var xxx_messageInfo_CurrencyRes proto.InternalMessageInfo

func (m *CurrencyRes) GetCode() uint32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *CurrencyRes) GetDesc() string {
	if m != nil {
		return m.Desc
	}
	return ""
}

type OnLineMsgReq struct {
	ToId                 []int64  `protobuf:"varint,1,rep,packed,name=to_id,json=toId" json:"to_id,omitempty"`
	MsgType              uint32   `protobuf:"varint,2,opt,name=MsgType" json:"MsgType,omitempty"`
	MsgContent           string   `protobuf:"bytes,3,opt,name=MsgContent" json:"MsgContent,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OnLineMsgReq) Reset()         { *m = OnLineMsgReq{} }
func (m *OnLineMsgReq) String() string { return proto.CompactTextString(m) }
func (*OnLineMsgReq) ProtoMessage()    {}
func (*OnLineMsgReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_api_698b2458fee38050, []int{4}
}
func (m *OnLineMsgReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OnLineMsgReq.Unmarshal(m, b)
}
func (m *OnLineMsgReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OnLineMsgReq.Marshal(b, m, deterministic)
}
func (dst *OnLineMsgReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OnLineMsgReq.Merge(dst, src)
}
func (m *OnLineMsgReq) XXX_Size() int {
	return xxx_messageInfo_OnLineMsgReq.Size(m)
}
func (m *OnLineMsgReq) XXX_DiscardUnknown() {
	xxx_messageInfo_OnLineMsgReq.DiscardUnknown(m)
}

var xxx_messageInfo_OnLineMsgReq proto.InternalMessageInfo

func (m *OnLineMsgReq) GetToId() []int64 {
	if m != nil {
		return m.ToId
	}
	return nil
}

func (m *OnLineMsgReq) GetMsgType() uint32 {
	if m != nil {
		return m.MsgType
	}
	return 0
}

func (m *OnLineMsgReq) GetMsgContent() string {
	if m != nil {
		return m.MsgContent
	}
	return ""
}

type OnLineMsgRes struct {
	SendUid              int64    `protobuf:"varint,1,opt,name=send_uid,json=sendUid" json:"send_uid,omitempty"`
	MsgType              uint32   `protobuf:"varint,2,opt,name=MsgType" json:"MsgType,omitempty"`
	MsgContent           string   `protobuf:"bytes,3,opt,name=MsgContent" json:"MsgContent,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OnLineMsgRes) Reset()         { *m = OnLineMsgRes{} }
func (m *OnLineMsgRes) String() string { return proto.CompactTextString(m) }
func (*OnLineMsgRes) ProtoMessage()    {}
func (*OnLineMsgRes) Descriptor() ([]byte, []int) {
	return fileDescriptor_api_698b2458fee38050, []int{5}
}
func (m *OnLineMsgRes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OnLineMsgRes.Unmarshal(m, b)
}
func (m *OnLineMsgRes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OnLineMsgRes.Marshal(b, m, deterministic)
}
func (dst *OnLineMsgRes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OnLineMsgRes.Merge(dst, src)
}
func (m *OnLineMsgRes) XXX_Size() int {
	return xxx_messageInfo_OnLineMsgRes.Size(m)
}
func (m *OnLineMsgRes) XXX_DiscardUnknown() {
	xxx_messageInfo_OnLineMsgRes.DiscardUnknown(m)
}

var xxx_messageInfo_OnLineMsgRes proto.InternalMessageInfo

func (m *OnLineMsgRes) GetSendUid() int64 {
	if m != nil {
		return m.SendUid
	}
	return 0
}

func (m *OnLineMsgRes) GetMsgType() uint32 {
	if m != nil {
		return m.MsgType
	}
	return 0
}

func (m *OnLineMsgRes) GetMsgContent() string {
	if m != nil {
		return m.MsgContent
	}
	return ""
}

type OffLineMsgRes struct {
	MsgType              uint32   `protobuf:"varint,1,opt,name=MsgType" json:"MsgType,omitempty"`
	MsgContent           string   `protobuf:"bytes,2,opt,name=MsgContent" json:"MsgContent,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OffLineMsgRes) Reset()         { *m = OffLineMsgRes{} }
func (m *OffLineMsgRes) String() string { return proto.CompactTextString(m) }
func (*OffLineMsgRes) ProtoMessage()    {}
func (*OffLineMsgRes) Descriptor() ([]byte, []int) {
	return fileDescriptor_api_698b2458fee38050, []int{6}
}
func (m *OffLineMsgRes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OffLineMsgRes.Unmarshal(m, b)
}
func (m *OffLineMsgRes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OffLineMsgRes.Marshal(b, m, deterministic)
}
func (dst *OffLineMsgRes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OffLineMsgRes.Merge(dst, src)
}
func (m *OffLineMsgRes) XXX_Size() int {
	return xxx_messageInfo_OffLineMsgRes.Size(m)
}
func (m *OffLineMsgRes) XXX_DiscardUnknown() {
	xxx_messageInfo_OffLineMsgRes.DiscardUnknown(m)
}

var xxx_messageInfo_OffLineMsgRes proto.InternalMessageInfo

func (m *OffLineMsgRes) GetMsgType() uint32 {
	if m != nil {
		return m.MsgType
	}
	return 0
}

func (m *OffLineMsgRes) GetMsgContent() string {
	if m != nil {
		return m.MsgContent
	}
	return ""
}

func init() {
	proto.RegisterType((*HeartbeatReq)(nil), "protos.HeartbeatReq")
	proto.RegisterType((*HeartbeatRes)(nil), "protos.HeartbeatRes")
	proto.RegisterType((*RegisterReq)(nil), "protos.RegisterReq")
	proto.RegisterType((*CurrencyRes)(nil), "protos.CurrencyRes")
	proto.RegisterType((*OnLineMsgReq)(nil), "protos.OnLineMsgReq")
	proto.RegisterType((*OnLineMsgRes)(nil), "protos.OnLineMsgRes")
	proto.RegisterType((*OffLineMsgRes)(nil), "protos.OffLineMsgRes")
	proto.RegisterEnum("protos.Tryout", Tryout_name, Tryout_value)
}

func init() { proto.RegisterFile("api.proto", fileDescriptor_api_698b2458fee38050) }

var fileDescriptor_api_698b2458fee38050 = []byte{
	// 315 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x92, 0xc1, 0x4e, 0xf2, 0x40,
	0x10, 0xc7, 0xbf, 0x52, 0x3e, 0x0a, 0x03, 0x08, 0x19, 0x30, 0xd4, 0x8b, 0x21, 0x3d, 0x18, 0xe2,
	0xc1, 0x8b, 0xf1, 0x09, 0x88, 0x89, 0x24, 0x12, 0x92, 0x0d, 0x1e, 0x0d, 0x81, 0xee, 0x50, 0x7b,
	0xd9, 0x85, 0x9d, 0xed, 0x81, 0x67, 0xf1, 0x45, 0x7c, 0x25, 0x7d, 0x0a, 0xd3, 0x6d, 0x34, 0x6b,
	0x0f, 0x5e, 0x3c, 0x75, 0xfe, 0xbf, 0x74, 0x7e, 0xb3, 0xd9, 0x59, 0xe8, 0x6c, 0x0f, 0xf9, 0xcd,
	0xc1, 0x68, 0xab, 0xb1, 0xe5, 0x3e, 0x9c, 0x9c, 0x41, 0xef, 0x81, 0xb6, 0xc6, 0xee, 0x68, 0x6b,
	0x05, 0x1d, 0x6b, 0x99, 0x93, 0x2b, 0xe8, 0x0a, 0xca, 0x72, 0xb6, 0x64, 0x04, 0x1d, 0x71, 0x02,
	0x51, 0xc1, 0x64, 0x36, 0xb9, 0x8c, 0x83, 0x69, 0x30, 0x0b, 0x45, 0xab, 0x8c, 0x0b, 0x99, 0xdc,
	0x41, 0x77, 0x5e, 0x18, 0x43, 0x2a, 0x3d, 0x09, 0x62, 0x44, 0x68, 0xa6, 0x5a, 0x92, 0xfb, 0xa9,
	0x2f, 0x5c, 0x5d, 0x32, 0x49, 0x9c, 0xc6, 0x8d, 0x69, 0x30, 0xeb, 0x08, 0x57, 0x27, 0xcf, 0xd0,
	0x5b, 0xa9, 0xc7, 0x5c, 0xd1, 0x92, 0xb3, 0xd2, 0x3f, 0x82, 0xff, 0x56, 0x57, 0xf6, 0x70, 0x16,
	0x8a, 0xa6, 0xd5, 0x0b, 0x89, 0x31, 0x44, 0x4b, 0xce, 0xd6, 0xa7, 0x03, 0xb9, 0xde, 0xbe, 0xf8,
	0x8a, 0x78, 0x09, 0xb0, 0xe4, 0x6c, 0xae, 0x95, 0x25, 0x65, 0xe3, 0xd0, 0x89, 0x3d, 0x92, 0xa4,
	0x3f, 0xf4, 0x8c, 0x17, 0xd0, 0x66, 0x52, 0x72, 0x53, 0x7c, 0x9f, 0x3f, 0x2a, 0xf3, 0x53, 0xfe,
	0x97, 0x21, 0x0b, 0xe8, 0xaf, 0xf6, 0x7b, 0x6f, 0x8a, 0xa7, 0x0a, 0x7e, 0x53, 0x35, 0xea, 0xaa,
	0xeb, 0xd7, 0x00, 0x5a, 0x6b, 0x73, 0xd2, 0x85, 0xc5, 0x1e, 0xb4, 0xef, 0xd3, 0x17, 0x5d, 0xb6,
	0x0d, 0xff, 0xe1, 0x18, 0x86, 0xfe, 0x9a, 0x1c, 0xa5, 0x1a, 0x65, 0x47, 0xf7, 0x38, 0x82, 0x81,
	0xb7, 0x32, 0x07, 0x33, 0x3c, 0x87, 0xa1, 0x7f, 0xd1, 0x8e, 0xbe, 0x47, 0x35, 0x5c, 0x19, 0x3e,
	0x22, 0x1c, 0xc3, 0xc0, 0xdb, 0xa6, 0xa3, 0x6f, 0x93, 0x5d, 0xf5, 0x66, 0x6e, 0x3f, 0x03, 0x00,
	0x00, 0xff, 0xff, 0xb8, 0x5a, 0xe2, 0xb2, 0x47, 0x02, 0x00, 0x00,
}
