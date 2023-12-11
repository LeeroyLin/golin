package net

type Message struct {
	MsgId   uint16
	ProtoId uint16
	MsgLen  uint16
	Data    []byte
}

func (m *Message) GetMsgId() uint16 {
	return m.MsgId
}

func (m *Message) SetMsgId(id uint16) {
	m.MsgId = id
}

func (m *Message) GetProtoId() uint16 {
	return m.ProtoId
}

func (m *Message) SetProtoId(id uint16) {
	m.ProtoId = id
}

func (m *Message) GetMsgLen() uint16 {
	return m.MsgLen
}

func (m *Message) SetMsgLen(len uint16) {
	m.MsgLen = len
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}
