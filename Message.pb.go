// Code generated by protoc-gen-go. DO NOT EDIT.
// source: Message.proto

package main

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type WSMessage struct {
	//消息类型 shop_auth ，auth，或 message
	Type                 string   `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Content              []byte   `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WSMessage) Reset()         { *m = WSMessage{} }
func (m *WSMessage) String() string { return proto.CompactTextString(m) }
func (*WSMessage) ProtoMessage()    {}
func (*WSMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_64b0f1bc979aed9f, []int{0}
}

func (m *WSMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WSMessage.Unmarshal(m, b)
}
func (m *WSMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WSMessage.Marshal(b, m, deterministic)
}
func (m *WSMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WSMessage.Merge(m, src)
}
func (m *WSMessage) XXX_Size() int {
	return xxx_messageInfo_WSMessage.Size(m)
}
func (m *WSMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_WSMessage.DiscardUnknown(m)
}

var xxx_messageInfo_WSMessage proto.InternalMessageInfo

func (m *WSMessage) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *WSMessage) GetContent() []byte {
	if m != nil {
		return m.Content
	}
	return nil
}

type AuthMessage struct {
	Token                string   `protobuf:"bytes,3,opt,name=token,proto3" json:"token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AuthMessage) Reset()         { *m = AuthMessage{} }
func (m *AuthMessage) String() string { return proto.CompactTextString(m) }
func (*AuthMessage) ProtoMessage()    {}
func (*AuthMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_64b0f1bc979aed9f, []int{1}
}

func (m *AuthMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AuthMessage.Unmarshal(m, b)
}
func (m *AuthMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AuthMessage.Marshal(b, m, deterministic)
}
func (m *AuthMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AuthMessage.Merge(m, src)
}
func (m *AuthMessage) XXX_Size() int {
	return xxx_messageInfo_AuthMessage.Size(m)
}
func (m *AuthMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_AuthMessage.DiscardUnknown(m)
}

var xxx_messageInfo_AuthMessage proto.InternalMessageInfo

func (m *AuthMessage) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type Message struct {
	//消息类型来源，buyer 或者 shoper
	From string `protobuf:"bytes,4,opt,name=from,proto3" json:"from,omitempty"`
	//消息类型，text 或者 img ,文字或图像
	MessageType          string   `protobuf:"bytes,5,opt,name=message_type,json=messageType,proto3" json:"message_type,omitempty"`
	FromUserId           int32    `protobuf:"varint,6,opt,name=from_user_id,json=fromUserId,proto3" json:"from_user_id,omitempty"`
	ToUserId             int32    `protobuf:"varint,7,opt,name=to_user_id,json=toUserId,proto3" json:"to_user_id,omitempty"`
	Content              string   `protobuf:"bytes,8,opt,name=content,proto3" json:"content,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_64b0f1bc979aed9f, []int{2}
}

func (m *Message) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message.Unmarshal(m, b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(b, m, deterministic)
}
func (m *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(m, src)
}
func (m *Message) XXX_Size() int {
	return xxx_messageInfo_Message.Size(m)
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetFrom() string {
	if m != nil {
		return m.From
	}
	return ""
}

func (m *Message) GetMessageType() string {
	if m != nil {
		return m.MessageType
	}
	return ""
}

func (m *Message) GetFromUserId() int32 {
	if m != nil {
		return m.FromUserId
	}
	return 0
}

func (m *Message) GetToUserId() int32 {
	if m != nil {
		return m.ToUserId
	}
	return 0
}

func (m *Message) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

func init() {
	proto.RegisterType((*WSMessage)(nil), "main.WSMessage")
	proto.RegisterType((*AuthMessage)(nil), "main.AuthMessage")
	proto.RegisterType((*Message)(nil), "main.Message")
}

func init() { proto.RegisterFile("Message.proto", fileDescriptor_64b0f1bc979aed9f) }

var fileDescriptor_64b0f1bc979aed9f = []byte{
	// 207 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xf5, 0x4d, 0x2d, 0x2e,
	0x4e, 0x4c, 0x4f, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xc9, 0x4d, 0xcc, 0xcc, 0x53,
	0xb2, 0xe4, 0xe2, 0x0c, 0x0f, 0x86, 0x4a, 0x08, 0x09, 0x71, 0xb1, 0x94, 0x54, 0x16, 0xa4, 0x4a,
	0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0x81, 0xd9, 0x42, 0x12, 0x5c, 0xec, 0xc9, 0xf9, 0x79, 0x25,
	0xa9, 0x79, 0x25, 0x12, 0x4c, 0x0a, 0x8c, 0x1a, 0x3c, 0x41, 0x30, 0xae, 0x92, 0x32, 0x17, 0xb7,
	0x63, 0x69, 0x49, 0x06, 0x4c, 0xb3, 0x08, 0x17, 0x6b, 0x49, 0x7e, 0x76, 0x6a, 0x9e, 0x04, 0x33,
	0x58, 0x37, 0x84, 0xa3, 0x34, 0x8b, 0x91, 0x8b, 0x1d, 0xc9, 0xf8, 0xb4, 0xa2, 0xfc, 0x5c, 0x09,
	0x16, 0x88, 0xf1, 0x20, 0xb6, 0x90, 0x22, 0x17, 0x4f, 0x2e, 0x44, 0x3a, 0x1e, 0x6c, 0x35, 0x2b,
	0x58, 0x8e, 0x1b, 0x2a, 0x16, 0x02, 0x72, 0x81, 0x02, 0x17, 0x0f, 0x48, 0x69, 0x7c, 0x69, 0x71,
	0x6a, 0x51, 0x7c, 0x66, 0x8a, 0x04, 0x9b, 0x02, 0xa3, 0x06, 0x6b, 0x10, 0x17, 0x48, 0x2c, 0xb4,
	0x38, 0xb5, 0xc8, 0x33, 0x45, 0x48, 0x86, 0x8b, 0xab, 0x24, 0x1f, 0x2e, 0xcf, 0x0e, 0x96, 0xe7,
	0x28, 0xc9, 0x87, 0xca, 0x22, 0xf9, 0x80, 0x03, 0x6c, 0x3a, 0x8c, 0x9b, 0xc4, 0x06, 0x0e, 0x09,
	0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0xfd, 0x87, 0xb3, 0xe7, 0x1a, 0x01, 0x00, 0x00,
}