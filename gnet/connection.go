package gnet

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/LeeroyLin/golin/iface"
	"github.com/LeeroyLin/golin/proto_model"
	"google.golang.org/protobuf/proto"
	"net"
)

type Connection struct {
	Conn      *net.TCPConn
	ConnId    uint32
	isClosed  bool
	ExitChan  chan bool
	RouterMap map[uint16]iface.IRouter
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

		go c.RouterHandle(&req)
	}
}

func (c *Connection) RouterHandle(request iface.IRequest) {
	msg := request.GetMsg()

	router, has := c.RouterMap[msg.GetProtoId()]
	if !has {
		return
	}

	fmt.Printf("【Req msg】 ConnId:%d MsgId:%d ProtoId:%d MsgLen:%d\n",
		c.GetConnId(),
		msg.GetMsgId(),
		msg.GetProtoId(),
		msg.GetMsgLen())

	jsonStr, err := json.Marshal(msg.GetData())
	if err != nil {
		fmt.Printf("Msg data to json failed.")
		return
	}

	fmt.Printf("【Req data】%s\n", jsonStr)

	router.PreHandle(request)

	errorCode, resData := router.Handle(request)

	if resData != nil {
		c.SendPB(request, errorCode, resData)
	}

	router.PostHandle(request)
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

	_, err = c.Conn.Write(nData)
	if err != nil {
		return errors.New("send msg error")
	}

	return nil
}

func (c *Connection) SendPB(request iface.IRequest, errorCode int32, resData proto.Message) error {
	pbResData, err := proto.Marshal(resData)
	if err != nil {
		fmt.Printf("Marshal protobuf data failed.")
		return err
	}

	respose := &proto_model.Response{
		Code: errorCode,
		Msg:  "",
		Data: pbResData,
	}

	pbResponse, err := proto.Marshal(respose)
	if err != nil {
		fmt.Printf("Marshal protobuf response failed.")
		return err
	}

	c.Send(request.GetMsg().GetProtoId()+1, pbResponse)

	return nil
}

func (c *Connection) SendPBNotify(resData proto.Message, errorCode int32) {
	pbResData, err := proto.Marshal(resData)
	if err != nil {
		fmt.Printf("Marshal protobuf data failed.")
		return
	}

	respose := &proto_model.Response{
		Code: errorCode,
		Msg:  "",
		Data: pbResData,
	}

	pbResponse, err := proto.Marshal(respose)
	if err != nil {
		fmt.Printf("Marshal protobuf response failed.")
		return
	}

	c.Send(0, pbResponse)
}

func NewConnection(conn *net.TCPConn, connId uint32, routerMap map[uint16]iface.IRouter) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnId:    connId,
		RouterMap: routerMap,
		isClosed:  false,
		ExitChan:  make(chan bool, 1),
	}

	return c
}

func (c *Connection) logfln(str string, a ...any) {
	fmt.Printf("[Conn:%d] %s\n", c.ConnId, fmt.Sprintf(str, a...))
}

func (c *Connection) logln(a ...any) {
	fmt.Printf("[Conn:%d] %s\n", c.ConnId, fmt.Sprint(a...))
}
