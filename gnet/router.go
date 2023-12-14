package gnet

import "github.com/LeeroyLin/golin/iface"

type BaseRouter struct {
}

func (br *BaseRouter) PreHandle(request iface.IRequest) {
}

func (br *BaseRouter) Handle(request iface.IRequest) {
}

func (br *BaseRouter) PostHandle(request iface.IRequest) {
}
