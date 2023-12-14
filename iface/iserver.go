package iface

type IServer interface {
	Start()
	Stop()
	Serve()
	AddRouter(protoId uint16, router IRouter)
}
