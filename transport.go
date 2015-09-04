package libnet

import (
	"io"
	"net"
	_ "time"
)

//transport 接口
type Transport interface {
	io.Reader
	io.Writer
	RemoteAddr() net.Addr
	Close() error
	Id() string
	ConnectTime() int64
	LastTime() int64
	Conn() net.Conn
	SetReadTimeout() error
	UnsetReadTimeout() error
	SetWriteTimeout() error
	UnsetWriteTimeout() error
}
