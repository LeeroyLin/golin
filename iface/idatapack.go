package iface

import "bytes"

type IDataPack interface {
	GetHeadLen() uint32
	Pack(msg IMessage) ([]byte, error)
	Unpack(*bytes.Buffer, []byte, *IConnection) (IMessage, error)
}
