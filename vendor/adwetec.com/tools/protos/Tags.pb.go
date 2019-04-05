// Code generated by protoc-gen-go. DO NOT EDIT.
// source: Tags.proto

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

type Tag struct {
	Tag                  *int64   `protobuf:"varint,1,req,name=tag" json:"tag,omitempty"`
	Time                 []int64  `protobuf:"varint,2,rep,name=time" json:"time,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Tag) Reset()         { *m = Tag{} }
func (m *Tag) String() string { return proto.CompactTextString(m) }
func (*Tag) ProtoMessage()    {}
func (*Tag) Descriptor() ([]byte, []int) {
	return fileDescriptor_Tags_ac32c08ba22d7323, []int{0}
}
func (m *Tag) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Tag.Unmarshal(m, b)
}
func (m *Tag) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Tag.Marshal(b, m, deterministic)
}
func (dst *Tag) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Tag.Merge(dst, src)
}
func (m *Tag) XXX_Size() int {
	return xxx_messageInfo_Tag.Size(m)
}
func (m *Tag) XXX_DiscardUnknown() {
	xxx_messageInfo_Tag.DiscardUnknown(m)
}

var xxx_messageInfo_Tag proto.InternalMessageInfo

func (m *Tag) GetTag() int64 {
	if m != nil && m.Tag != nil {
		return *m.Tag
	}
	return 0
}

func (m *Tag) GetTime() []int64 {
	if m != nil {
		return m.Time
	}
	return nil
}

type Tags struct {
	Tags                 []*Tag   `protobuf:"bytes,1,rep,name=tags" json:"tags,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Tags) Reset()         { *m = Tags{} }
func (m *Tags) String() string { return proto.CompactTextString(m) }
func (*Tags) ProtoMessage()    {}
func (*Tags) Descriptor() ([]byte, []int) {
	return fileDescriptor_Tags_ac32c08ba22d7323, []int{1}
}
func (m *Tags) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Tags.Unmarshal(m, b)
}
func (m *Tags) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Tags.Marshal(b, m, deterministic)
}
func (dst *Tags) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Tags.Merge(dst, src)
}
func (m *Tags) XXX_Size() int {
	return xxx_messageInfo_Tags.Size(m)
}
func (m *Tags) XXX_DiscardUnknown() {
	xxx_messageInfo_Tags.DiscardUnknown(m)
}

var xxx_messageInfo_Tags proto.InternalMessageInfo

func (m *Tags) GetTags() []*Tag {
	if m != nil {
		return m.Tags
	}
	return nil
}

func init() {
	proto.RegisterType((*Tag)(nil), "protos.Tag")
	proto.RegisterType((*Tags)(nil), "protos.Tags")
}

func init() { proto.RegisterFile("Tags.proto", fileDescriptor_Tags_ac32c08ba22d7323) }

var fileDescriptor_Tags_ac32c08ba22d7323 = []byte{
	// 104 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x0a, 0x49, 0x4c, 0x2f,
	0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x03, 0x53, 0xc5, 0x4a, 0x0a, 0x5c, 0xcc, 0x21,
	0x89, 0xe9, 0x42, 0xdc, 0x5c, 0xcc, 0x25, 0x89, 0xe9, 0x12, 0x8c, 0x0a, 0x4c, 0x1a, 0xcc, 0x42,
	0x3c, 0x5c, 0x2c, 0x25, 0x99, 0xb9, 0xa9, 0x12, 0x4c, 0x0a, 0xcc, 0x1a, 0xcc, 0x4a, 0x8a, 0x5c,
	0x2c, 0x20, 0x7d, 0x42, 0x92, 0x5c, 0x2c, 0x25, 0x89, 0xe9, 0xc5, 0x12, 0x8c, 0x0a, 0xcc, 0x1a,
	0xdc, 0x46, 0xdc, 0x10, 0x73, 0x8a, 0xf5, 0x42, 0x12, 0xd3, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff,
	0x07, 0x46, 0x66, 0x2e, 0x59, 0x00, 0x00, 0x00,
}