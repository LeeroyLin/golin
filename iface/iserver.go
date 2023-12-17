package iface

import "google.golang.org/protobuf/proto"

type IServer interface {
	Start()
	Stop()
	Serve()
	AddRouter(protoId uint16, router IRouter, reqData proto.Message)
	GetConnMgr() IConnManager
	SetOnConnStart(handler func(conn IConnection))
	CallOnConnStart(conn IConnection)
	SetOnConnStop(handler func(conn IConnection))
	CallOnConnStop(conn IConnection)
}
