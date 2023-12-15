package gnet

import (
	"encoding/json"
	"github.com/LeeroyLin/golin/iface"
	"google.golang.org/protobuf/proto"
)

type MessageHandler struct {
	Apis map[uint16]RouterData
}

func (h *MessageHandler) RouterHandle(request iface.IRequest) {
	c := request.GetConnection()
	msg := request.GetMsg()

	c.Logfln("【Req msg】 ConnId:%d MsgId:%d ProtoId:%d MsgLen:%d",
		c.GetConnId(),
		msg.GetMsgId(),
		msg.GetProtoId(),
		msg.GetMsgLen())

	routerData, has := h.Apis[msg.GetProtoId()]
	if !has {
		return
	}

	err := proto.Unmarshal(msg.GetData(), routerData.ReqData)
	if err != nil {
		c.Logln("Msg data to protobuf data failed. err: ", err)
		return
	}

	jsonStr, err := json.Marshal(routerData.ReqData)
	if err != nil {
		c.Logln("Msg data to json failed.")
		return
	}

	c.Logfln("【Req data】%s", jsonStr)

	routerData.Router.PreHandle(request, routerData.ReqData)

	errorCode, resData, resMsg := routerData.Router.Handle(request, routerData.ReqData)

	err = c.SendPB(request, errorCode, resData, resMsg)
	if err != nil {
		c.Logfln("Send pb failed. err: ", err)
		return
	}

	routerData.Router.PostHandle(request, routerData.ReqData)
}

func (h *MessageHandler) Has(protoId uint16) bool {
	_, has := h.Apis[protoId]
	return has
}

func (h *MessageHandler) Add(protoId uint16, router iface.IRouter, reqData proto.Message) {
	h.Apis[protoId] = RouterData{
		Router:  router,
		ReqData: reqData,
	}
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{
		Apis: make(map[uint16]RouterData),
	}
}
