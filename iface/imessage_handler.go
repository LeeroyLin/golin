package iface

import "google.golang.org/protobuf/proto"

type IMessageHandler interface {
	DoMsgHandle(request IRequest)
	Has(protoId uint16) bool
	Add(protoId uint16, router IRouter, reqData proto.Message)
	StartWorkerPool()
}
