package iface

import "google.golang.org/protobuf/proto"

type IRouter interface {
	PreHandle(request IRequest)
	Handle(request IRequest) (int32, proto.Message)
	PostHandle(request IRequest)
}
