package iface

import "google.golang.org/protobuf/proto"

type IServer interface {
	Start()
	Stop()
	Serve()
	AddRouter(protoId uint16, router IRouter, reqData proto.Message)
	GetConnMgr() IConnManager
}
