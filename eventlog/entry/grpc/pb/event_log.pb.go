// Copyright (C) 2021 Gridworkz Co., Ltd.
// KATO, Application Management Platform

// Permission is hereby granted, free of charge, to any person obtaining a copy of this 
// software and associated documentation files (the "Software"), to deal in the Software
// without restriction, including without limitation the rights to use, copy, modify, merge,
// publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons 
// to whom the Software is furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all copies or 
// substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, 
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
// PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE
// FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type LogMessage struct {
	Log []byte `protobuf:"bytes,1,opt,name=log,proto3" json:"log,omitempty"`
}

func (m *LogMessage) Reset()                    { *m = LogMessage{} }
func (m *LogMessage) String() string            { return proto.CompactTextString(m) }
func (*LogMessage) ProtoMessage()               {}
func (*LogMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *LogMessage) GetLog() []byte {
	if m != nil {
		return m.Log
	}
	return nil
}

type Reply struct {
	Status  string `protobuf:"bytes,1,opt,name=status" json:"status,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message" json:"message,omitempty"`
}

func (m *Reply) Reset()                    { *m = Reply{} }
func (m *Reply) String() string            { return proto.CompactTextString(m) }
func (*Reply) ProtoMessage()               {}
func (*Reply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Reply) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *Reply) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*LogMessage)(nil), "pb.LogMessage")
	proto.RegisterType((*Reply)(nil), "pb.Reply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for EventLog service

type EventLogClient interface {
	Log(ctx context.Context, opts ...grpc.CallOption) (EventLog_LogClient, error)
}

type eventLogClient struct {
	cc *grpc.ClientConn
}

func NewEventLogClient(cc *grpc.ClientConn) EventLogClient {
	return &eventLogClient{cc}
}

func (c *eventLogClient) Log(ctx context.Context, opts ...grpc.CallOption) (EventLog_LogClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_EventLog_serviceDesc.Streams[0], c.cc, "/pb.EventLog/Log", opts...)
	if err != nil {
		return nil, err
	}
	x := &eventLogLogClient{stream}
	return x, nil
}

type EventLog_LogClient interface {
	Send(*LogMessage) error
	CloseAndRecv() (*Reply, error)
	grpc.ClientStream
}

type eventLogLogClient struct {
	grpc.ClientStream
}

func (x *eventLogLogClient) Send(m *LogMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *eventLogLogClient) CloseAndRecv() (*Reply, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Reply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for EventLog service

type EventLogServer interface {
	Log(EventLog_LogServer) error
}

func RegisterEventLogServer(s *grpc.Server, srv EventLogServer) {
	s.RegisterService(&_EventLog_serviceDesc, srv)
}

func _EventLog_Log_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(EventLogServer).Log(&eventLogLogServer{stream})
}

type EventLog_LogServer interface {
	SendAndClose(*Reply) error
	Recv() (*LogMessage, error)
	grpc.ServerStream
}

type eventLogLogServer struct {
	grpc.ServerStream
}

func (x *eventLogLogServer) SendAndClose(m *Reply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *eventLogLogServer) Recv() (*LogMessage, error) {
	m := new(LogMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _EventLog_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.EventLog",
	HandlerType: (*EventLogServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Log",
			Handler:       _EventLog_Log_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "event_log.proto",
}

func init() { proto.RegisterFile("event_log.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 150 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4f, 0x2d, 0x4b, 0xcd,
	0x2b, 0x89, 0xcf, 0xc9, 0x4f, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x48, 0x52,
	0x92, 0xe3, 0xe2, 0xf2, 0xc9, 0x4f, 0xf7, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0x15, 0x12, 0xe0,
	0x62, 0xce, 0xc9, 0x4f, 0x97, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x09, 0x02, 0x31, 0x95, 0x2c, 0xb9,
	0x58, 0x83, 0x52, 0x0b, 0x72, 0x2a, 0x85, 0xc4, 0xb8, 0xd8, 0x8a, 0x4b, 0x12, 0x4b, 0x4a, 0x8b,
	0xc1, 0xb2, 0x9c, 0x41, 0x50, 0x9e, 0x90, 0x04, 0x17, 0x7b, 0x2e, 0x44, 0xb7, 0x04, 0x13, 0x58,
	0x02, 0xc6, 0x35, 0x32, 0xe0, 0xe2, 0x70, 0x05, 0xd9, 0xe8, 0x93, 0x9f, 0x2e, 0xa4, 0xc2, 0xc5,
	0x0c, 0xa2, 0xf8, 0xf4, 0x0a, 0x92, 0xf4, 0x10, 0xf6, 0x49, 0x71, 0x82, 0xf8, 0x60, 0xf3, 0x95,
	0x18, 0x34, 0x18, 0x93, 0xd8, 0xc0, 0xee, 0x32, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0xd5, 0xa6,
	0xad, 0x69, 0xaa, 0x00, 0x00, 0x00,
}
