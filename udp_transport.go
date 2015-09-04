package libnet

import (
	"code.google.com/p/go-uuid/uuid"
	"net"
	"time"
)

type UdpTransport struct {
	conn *net.UDPConn
	*BaseTransport
}

func NewUDPTransport(conn *net.UDPConn) *UdpTransport {
	trans := &UdpTransport{}
	trans.conn = conn

	trans.BaseTransport = new(BaseTransport)
	trans.id = uuid.New()
	trans.connectTime = time.Now().Unix()
	trans.lastTime = 0
	return trans
}

func (t *UdpTransport) String() string {
	return "id:" + t.id + " addr:" + t.conn.RemoteAddr().String()
}

//实现 io.Writer
func (t *UdpTransport) Write(b []byte) (int, error) {
	t.lastTime = time.Now().Unix()
	return t.conn.Write(b)
}

//实现 io.Reader
func (t *UdpTransport) Read(b []byte) (int, error) {
	cnt, err := t.conn.Read(b)
	return cnt, err
}

func (t *UdpTransport) RemoteAddr() net.Addr {
	return t.conn.RemoteAddr()
}

func (t *UdpTransport) Close() error {
	return t.conn.Close()
}

func (t *UdpTransport) Id() string {
	return t.id
}

func (t *UdpTransport) Conn() net.Conn {
	return t.conn
}

func (t *UdpTransport) SetReadTimeout() error {

	return nil
}

func (t *UdpTransport) UnsetReadTimeout() error {
	return nil
}

func (t *UdpTransport) SetWriteTimeout() error {
	return nil
}

func (t *UdpTransport) UnsetWriteTimeout() error {
	return nil
}

func (t *UdpTransport) ConnectTime() int64 {
	return t.connectTime
}

func (t *UdpTransport) LastTime() int64 {
	return t.lastTime
}
