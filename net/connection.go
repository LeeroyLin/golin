package net

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/LeeroyLin/golin/iface"
	"net"
)

type Connection struct {
	Conn     *net.TCPConn
	ConnId   uint32
	isClosed bool
	ExitChan chan bool
	Router   iface.IRouter
}

func (c *Connection) StartReader() {
	defer c.Stop()

	c.logfln("Start reader")

	binaryData := make([]byte, 1024)

	dataPack := DataPack{}

	buffer := bytes.NewBuffer([]byte{})

	targetHeadLen := dataPack.GetHeadLen()

	for {
		if buffer.Len() < targetHeadLen {
			cnt, err := c.Conn.Read(binaryData)
			if err != nil {
				c.logln("Read buf err:", err)
				return
			}

			if cnt == 0 {
				continue
			}

			buffer.Write(binaryData[:cnt])

			if buffer.Len() < targetHeadLen {
				c.logln("Msg header length %d less than target length %d", buffer.Len(), targetHeadLen)
				return
			}
		}

		msg, err := dataPack.Unpack(buffer, binaryData, c)
		if err != nil {
			c.logln("Unpack msg data err:", err)

			return
		}

		req := Request{
			conn: c,
			msg:  msg,
		}

		go func(request iface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
	}
}

func (c *Connection) Start() {
	c.logln("【Start】 connection")

	go c.StartReader()
}

func (c *Connection) Stop() {
	c.logln("【Stop】 connection")

	if c.isClosed {
		return
	}

	c.isClosed = true

	c.Conn.Close()

	close(c.ExitChan)
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
	dataPack := DataPack{}

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

	_, err = c.Conn.Write(nData)
	if err != nil {
		return errors.New("send msg error")
	}

	return nil
}

func NewConnection(conn *net.TCPConn, connId uint32, router iface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnId:   connId,
		Router:   router,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}

	return c
}

func (c *Connection) logfln(str string, a ...any) {
	fmt.Printf("[Conn:%d] %s\n", c.ConnId, fmt.Sprintf(str, a...))
}

func (c *Connection) logln(a ...any) {
	fmt.Printf("[Conn:%d] %s\n", c.ConnId, fmt.Sprint(a...))
}
