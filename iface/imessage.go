package iface

type IMessage interface {
	GetMsgId() uint16
	SetMsgId(id uint16)
	GetProtoId() uint16
	SetProtoId(id uint16)
	GetMsgLen() uint16
	SetMsgLen(len uint16)
	GetData() []byte
	SetData(data []byte)
}
