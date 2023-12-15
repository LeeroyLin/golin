package iface

import (
	"google.golang.org/protobuf/proto"
	"net"
)

type IConnection interface {
	Start()
	Stop()
	GetTCPConnection() *net.TCPConn
	GetConnId() uint32
	GetRemoteAddr() net.Addr
	Send(protoId uint16, data []byte) error
	SendPB(request IRequest, errorCode int32, resData proto.Message, resMsg string) error
}
