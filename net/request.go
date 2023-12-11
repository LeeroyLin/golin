package net

import "github.com/LeeroyLin/golin/iface"

type Request struct {
	conn iface.IConnection
	msg  iface.IMessage
}

func (r *Request) GetConnection() iface.IConnection {
	return r.conn
}

func (r *Request) GetMsg() iface.IMessage {
	return r.msg
}
