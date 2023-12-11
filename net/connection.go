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
	defer fmt.Println("Finish conn read. id=", c.ConnId, ", remove addr=", c.GetRemoteAddr().String())
	defer c.Stop()

	fmt.Println("Start conn read. id=", c.ConnId)

	binaryData := make([]byte, 1024)

	dataPack := DataPack{}

	buffer := bytes.NewBuffer([]byte{})

	for {
		if buffer.Len() < int(dataPack.GetHeadLen()) {
			cnt, err := c.Conn.Read(binaryData)
			if err != nil {
				fmt.Println("Read buf err", err, ", id=", c.ConnId)
				continue
			}

			if cnt == 0 {
				continue
			}

			buffer.Write(binaryData[:cnt])

			if buffer.Len() < int(dataPack.GetHeadLen()) {
				continue
			}
		}

		msg, err := dataPack.Unpack(buffer, binaryData, c)
		if err != nil {
			fmt.Println("Unpack msg data err", err, ", id=", c.ConnId)
			continue
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
	fmt.Println("Start conn id=", c.ConnId)

	go c.StartReader()
}

func (c *Connection) Stop() {
	fmt.Println("Stop conn id=", c.ConnId)

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

func (c *Connection) Send(data []byte) error {
	_, err := c.Conn.Write(data)
	if err != nil {
		return errors.New("conn send error")
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
