package gnet

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/LeeroyLin/golin/iface"
	"github.com/LeeroyLin/golin/proto_model"
	"google.golang.org/protobuf/proto"
	"net"
)

type Connection struct {
	Conn           *net.TCPConn
	ConnId         uint32
	isClosed       bool
	exitChan       chan bool
	msgChan        chan []byte
	MessageHandler iface.IMessageHandler
}

func (c *Connection) StartReader() {
	defer c.Stop()

	c.Logfln("Start reader")

	binaryData := make([]byte, 1024)

	dataPack := DataPack{}

	buffer := bytes.NewBuffer([]byte{})

	targetHeadLen := dataPack.GetHeadLen()

	for {
		if buffer.Len() < targetHeadLen {
			cnt, err := c.Conn.Read(binaryData)
			if err != nil {
				c.Logln("Read buf err:", err)
				return
			}

			if cnt == 0 {
				continue
			}

			buffer.Write(binaryData[:cnt])

			if buffer.Len() < targetHeadLen {
				c.Logln("Msg header length %d less than target length %d", buffer.Len(), targetHeadLen)
				return
			}
		}

		msg, err := dataPack.Unpack(buffer, binaryData, c)
		if err != nil {
			c.Logln("Unpack msg data err:", err)

			return
		}

		req := Request{
			conn: c,
			msg:  msg,
		}

		go c.MessageHandler.RouterHandle(&req)
	}
}

func (c *Connection) StartWriter() {
	c.Logfln("Start writer")

	for {
		select {
		case m := <-c.msgChan:
			_, err := c.Conn.Write(m)
			if err != nil {
				c.Logln("send msg failed. err: ", err)
			}
			break
		case <-c.exitChan:
			break
		}
	}
}

func (c *Connection) Start() {
	c.Logln("【Start】 connection")

	go c.StartReader()
}

func (c *Connection) Stop() {
	c.Logln("【Stop】 connection")

	if c.isClosed {
		return
	}

	c.isClosed = true

	err := c.Conn.Close()
	if err != nil {
		c.Logfln("Close connection failed. err: ", err)
		return
	}

	close(c.exitChan)
	close(c.msgChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnId() uint32 {
	return c.ConnId
}

func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(protoId uint16, data []byte) error {
	dataPack := NewDataPack()

	msg := &Message{
		MsgId:   0,
		ProtoId: protoId,
		MsgLen:  uint32(len(data)),
		Data:    data,
	}

	nData, err := dataPack.Pack(msg)
	if err != nil {
		return errors.New("pack msg error")
	}

	c.msgChan <- nData

	return nil
}

func (c *Connection) SendPB(request iface.IRequest, errorCode int32, resData proto.Message, resMsg string) error {
	var pbResData []byte
	if resData != nil {
		data, err := proto.Marshal(resData)
		if err != nil {
			c.Logln("Marshal protobuf data failed.")
			return err
		}

		pbResData = data
	}

	res := &proto_model.ProtoResponse{
		Code: errorCode,
		Msg:  resMsg,
		Data: pbResData,
	}

	pbResponse, err := proto.Marshal(res)
	if err != nil {
		c.Logln("Marshal protobuf response failed.")
		return err
	}

	return c.Send(request.GetMsg().GetProtoId()+1, pbResponse)
}

func (c *Connection) SendPBNotify(resData proto.Message, errorCode int32) error {
	pbResData, err := proto.Marshal(resData)
	if err != nil {
		c.Logln("Marshal protobuf data failed.")
		return err
	}

	res := &proto_model.ProtoResponse{
		Code: errorCode,
		Msg:  "",
		Data: pbResData,
	}

	pbResponse, err := proto.Marshal(res)
	if err != nil {
		c.Logln("Marshal protobuf response failed.")
		return err
	}

	return c.Send(0, pbResponse)
}

func NewConnection(conn *net.TCPConn, connId uint32, MessageHandler iface.IMessageHandler) *Connection {
	c := &Connection{
		Conn:           conn,
		ConnId:         connId,
		MessageHandler: MessageHandler,
		isClosed:       false,
		exitChan:       make(chan bool, 1),
		msgChan:        make(chan []byte),
	}

	return c
}

func (c *Connection) Logfln(str string, a ...any) {
	fmt.Printf("[Conn:%d] %s\n", c.ConnId, fmt.Sprintf(str, a...))
}

func (c *Connection) Logln(a ...any) {
	fmt.Printf("[Conn:%d] %s\n", c.ConnId, fmt.Sprint(a...))
}
