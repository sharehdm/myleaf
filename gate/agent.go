package gate

import (
	"net"
	"time"
)

type Agent interface {
	WriteMsg(pid uint16, mid uint16, msg interface{}) // 发送消息
	WriteBytes(msg ...[]byte)                         // 发送消息
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	Destroy()
	UserData() interface{}
	SetUserData(data interface{})
	LoginCheck(count time.Duration) // 登录检测，多长时间不登陆断开连接
	CancelLoginCheck()
}
