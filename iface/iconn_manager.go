package iface

type IConnManager interface {
	Add(conn IConnection)
	Remove(conn IConnection)
	Get(connId uint32) (conn IConnection, err error)
	Len() int
	ClearConn()
}
