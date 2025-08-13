package gate

import (
	"net"
)

type Agent interface {
	WriteMsg(pid uint16, mid uint16, msg interface{})
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	Destroy()
	UserData() interface{}
	SetUserData(data interface{})
}
