package iface

import "google.golang.org/protobuf/proto"

type IRouter interface {
	PreHandle(request IRequest, reqData interface{})
	Handle(request IRequest, reqData interface{}) (int32, proto.Message, string)
	PostHandle(request IRequest, reqData interface{})
}
