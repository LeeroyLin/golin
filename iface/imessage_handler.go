package iface

import "google.golang.org/protobuf/proto"

type IMessageHandler interface {
	RouterHandle(request IRequest)
	Has(protoId uint16) bool
	Add(protoId uint16, router IRouter, reqData proto.Message)
}
