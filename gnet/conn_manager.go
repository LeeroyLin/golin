package gnet

import (
	"errors"
	"fmt"
	"github.com/LeeroyLin/golin/iface"
	"sync"
)

type ConnManger struct {
	connections map[uint32]iface.IConnection
	connLock    sync.RWMutex
}

func (m *ConnManger) Add(conn iface.IConnection) {
	m.connLock.Lock()
	defer m.connLock.Unlock()

	m.connections[conn.GetConnId()] = conn

	conn.Logln("connection added. id: ", conn.GetConnId())
}

func (m *ConnManger) Remove(conn iface.IConnection) {
	m.connLock.Lock()
	defer m.connLock.Unlock()

	delete(m.connections, conn.GetConnId())

	conn.Logln("connection removed. id: ", conn.GetConnId())
}

func (m *ConnManger) Get(connId uint32) (conn iface.IConnection, err error) {
	m.connLock.RLock()
	defer m.connLock.RUnlock()

	conn, ok := m.connections[connId]
	if ok {
		return conn, nil
	}

	return nil, errors.New(fmt.Sprintf("Not find conn with id %d", connId))
}

func (m *ConnManger) Len() int {
	return len(m.connections)
}

func (m *ConnManger) ClearConn() {
	m.connLock.Lock()
	defer m.connLock.Unlock()

	for connId, conn := range m.connections {
		conn.Stop()
		delete(m.connections, connId)
	}

	fmt.Println("ConnManager clear connections.")
}

func NewConnManager() iface.IConnManager {
	return &ConnManger{
		connections: make(map[uint32]iface.IConnection),
	}
}
