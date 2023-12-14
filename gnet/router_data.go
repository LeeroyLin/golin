package gnet

import (
	"github.com/LeeroyLin/golin/iface"
	"google.golang.org/protobuf/proto"
)

type RouterData struct {
	Router  iface.IRouter
	ReqData proto.Message
}
