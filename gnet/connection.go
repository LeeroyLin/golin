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
	exitChan  chan bool
	msgChan   chan []byte
	RouterMap map[uint16]RouterData
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

func (c *Connection) StartWriter() {
	c.logfln("Start writer")

	for {
		select {
		case m := <-c.msgChan:
			_, err := c.Conn.Write(m)
			if err != nil {
				c.logln("send msg failed. err: ", err)
			}
			break
		case <-c.exitChan:
			break
		}
	}
}

func (c *Connection) RouterHandle(request iface.IRequest) {
	msg := request.GetMsg()

	c.logfln("【Req msg】 ConnId:%d MsgId:%d ProtoId:%d MsgLen:%d",
		c.GetConnId(),
		msg.GetMsgId(),
		msg.GetProtoId(),
		msg.GetMsgLen())

	routerData, has := c.RouterMap[msg.GetProtoId()]
	if !has {
		return
	}

	err := proto.Unmarshal(msg.GetData(), routerData.ReqData)
	if err != nil {
		c.logln("Msg data to protobuf data failed. err: ", err)
		return
	}

	jsonStr, err := json.Marshal(routerData.ReqData)
	if err != nil {
		c.logln("Msg data to json failed.")
		return
	}

	c.logfln("【Req data】%s", jsonStr)

	routerData.Router.PreHandle(request, routerData.ReqData)

	errorCode, resData, resMsg := routerData.Router.Handle(request, routerData.ReqData)

	err = c.SendPB(request, errorCode, resData, resMsg)
	if err != nil {
		c.logfln("Send pb failed. err: ", err)
		return
	}

	routerData.Router.PostHandle(request, routerData.ReqData)
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

	err := c.Conn.Close()
	if err != nil {
		c.logfln("Close connection failed. err: ", err)
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
			c.logln("Marshal protobuf data failed.")
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
		c.logln("Marshal protobuf response failed.")
		return err
	}

	return c.Send(request.GetMsg().GetProtoId()+1, pbResponse)
}

func (c *Connection) SendPBNotify(resData proto.Message, errorCode int32) error {
	pbResData, err := proto.Marshal(resData)
	if err != nil {
		c.logln("Marshal protobuf data failed.")
		return err
	}

	res := &proto_model.ProtoResponse{
		Code: errorCode,
		Msg:  "",
		Data: pbResData,
	}

	pbResponse, err := proto.Marshal(res)
	if err != nil {
		c.logln("Marshal protobuf response failed.")
		return err
	}

	return c.Send(0, pbResponse)
}

func NewConnection(conn *net.TCPConn, connId uint32, routerMap map[uint16]RouterData) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnId:    connId,
		RouterMap: routerMap,
		isClosed:  false,
		exitChan:  make(chan bool, 1),
		msgChan:   make(chan []byte),
	}

	return c
}

func (c *Connection) logfln(str string, a ...any) {
	fmt.Printf("[Conn:%d] %s\n", c.ConnId, fmt.Sprintf(str, a...))
}

func (c *Connection) logln(a ...any) {
	fmt.Printf("[Conn:%d] %s\n", c.ConnId, fmt.Sprint(a...))
}
