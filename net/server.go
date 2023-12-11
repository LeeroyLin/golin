package net

import (
	"fmt"
	"golin/iface"
	"golin/utils"
	"net"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
	Router    iface.IRouter
}

func (s *Server) AddRouter(router iface.IRouter) {
	s.Router = router
}

func (s *Server) Start() {
	fmt.Printf("Start %s server...\nHost:%s Port:%d\n",
		utils.GlobalConfig.Name, utils.GlobalConfig.Host, utils.GlobalConfig.TcpPort)

	fmt.Printf("MaxConn:%d\n", utils.GlobalConfig.MaxConn)

	go func() {
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("Receive tcp addr error: ", err)
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

			dealConn := NewConnection(conn, cid, s.Router)
			cid++

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
		Name:      utils.GlobalConfig.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalConfig.Host,
		Port:      utils.GlobalConfig.TcpPort,
		Router:    nil,
	}

	return s
}