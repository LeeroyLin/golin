package gnet

import (
	"fmt"
	"github.com/LeeroyLin/golin/iface"
	"github.com/LeeroyLin/golin/utils"
	"google.golang.org/protobuf/proto"
	"net"
)

type Server struct {
	Name           string
	IPVersion      string
	IP             string
	Port           int
	MessageHandler iface.IMessageHandler
}

func (s *Server) AddRouter(protoId uint16, router iface.IRouter, reqData proto.Message) {
	if s.MessageHandler.Has(protoId) {
		fmt.Println("Already has router handle at protoId: ", protoId)
		return
	}

	s.MessageHandler.Add(protoId, router, reqData)
}

func (s *Server) Start() {
	fmt.Printf("Start '%s' server...\nHost:%s Port:%d\n",
		utils.GlobalConfig.Name, utils.GlobalConfig.Host, utils.GlobalConfig.TcpPort)

	fmt.Printf("MaxConn:%d\n", utils.GlobalConfig.MaxConn)
	fmt.Printf("MaxMsgLen:%d\n", utils.GlobalConfig.MaxMsgLen)
	fmt.Printf("IsEncrypt:%t\n", utils.GlobalConfig.IsEncrypt)
	fmt.Printf("RC4Key:%s\n", utils.GlobalConfig.RC4Key)
	fmt.Printf("WorkerPoolSize:%d\n", utils.GlobalConfig.WorkerPoolSize)

	go func() {
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("Resolve tcp addr error: ", err)
			return
		}

		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("Listen ", s.IP, ":", s.Port, " error: ", err)
			return
		}

		fmt.Printf("Listen to %s:%d success\n", s.IP, s.Port)

		var cid uint32
		cid = 0

		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept tcp error", err)

				continue
			}

			dealConn := NewConnection(conn, cid, s.MessageHandler)

			if cid == ^uint32(0) {
				cid = 0
			} else {
				cid++
			}

			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
}

func (s *Server) Serve() {
	s.Start()

	select {}
}

func NewServer() iface.IServer {
	s := &Server{
		Name:           utils.GlobalConfig.Name,
		IPVersion:      "tcp4",
		IP:             utils.GlobalConfig.Host,
		Port:           utils.GlobalConfig.TcpPort,
		MessageHandler: NewMessageHandler(),
	}

	return s
}
