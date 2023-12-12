package iface

type IMessage interface {
	GetMsgId() uint16
	SetMsgId(id uint16)
	GetProtoId() uint16
	SetProtoId(id uint16)
	GetMsgLen() uint32
	SetMsgLen(len uint32)
	GetData() []byte
	SetData(data []byte)
}
