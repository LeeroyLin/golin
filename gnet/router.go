package gnet

import (
	"github.com/LeeroyLin/golin/iface"
	"google.golang.org/protobuf/proto"
)

type BaseRouter struct {
}

func (br *BaseRouter) PreHandle(request iface.IRequest, reqData interface{}) {
}

func (br *BaseRouter) Handle(request iface.IRequest, reqData interface{}) (int32, proto.Message) {
	return 0, nil
}

func (br *BaseRouter) PostHandle(request iface.IRequest, reqData interface{}) {
}
