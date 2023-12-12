package net

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/LeeroyLin/golin/iface"
	"github.com/LeeroyLin/golin/utils"
)

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() int {
	// 序列号2字节 协议号2字节 内容长度4字节
	return 8
}

func (dp *DataPack) Pack(msg iface.IMessage) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})

	// 写 序列号
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	// 写 协议号
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetProtoId()); err != nil {
		return nil, err
	}

	// 写 内容长度
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}

	// 写 内容
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

func (dp *DataPack) Unpack(dataBuffer *bytes.Buffer, binaryData []byte, c *Connection) (iface.IMessage, error) {
	msg := &Message{}

	// 读 序列号
	if err := binary.Read(dataBuffer, binary.LittleEndian, &msg.MsgId); err != nil {
		return nil, err
	}

	// 读 协议号
	if err := binary.Read(dataBuffer, binary.LittleEndian, &msg.ProtoId); err != nil {
		return nil, err
	}

	// 读 内容长度
	if err := binary.Read(dataBuffer, binary.LittleEndian, &msg.MsgLen); err != nil {
		return nil, err
	}

	// 判断消息内容长度
	if msg.MsgLen > utils.GlobalConfig.MaxMsgLen {
		return nil, errors.New(fmt.Sprintf(
			"Msg length %d over max length %d.",
			msg.MsgLen,
			utils.GlobalConfig.MaxMsgLen,
		))
	}

	msg.SetData(make([]byte, msg.MsgLen))

	for {
		if dataBuffer.Len() >= int(msg.MsgLen) {
			break
		}

		cnt, err := c.Conn.Read(binaryData)
		if err != nil {
			return nil, errors.New(fmt.Sprint("Read msg content err:", err))
		}

		if cnt == 0 {
			return nil, errors.New(fmt.Sprint("Read msg content failed. no enough data."))
		}

		dataBuffer.Write(binaryData[:cnt])
	}

	// 读 内容
	if err := binary.Read(dataBuffer, binary.LittleEndian, &msg.Data); err != nil {
		return nil, err
	}

	return msg, nil
}
