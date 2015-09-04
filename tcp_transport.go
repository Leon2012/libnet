package libnet

import (
	"code.google.com/p/go-uuid/uuid"
	"net"
	"time"
)

type TcpTransport struct {
	conn              *net.TCPConn
	hasSetReadTimeout bool
	readTimeout       time.Duration
	writeTimeout      time.Duration
	*BaseTransport
}

func NewTCPTransport(conn *net.TCPConn, rTimeout time.Duration, wTimeout time.Duration) *TcpTransport {
	trans := &TcpTransport{}
	trans.hasSetReadTimeout = false
	trans.readTimeout = rTimeout
	trans.writeTimeout = wTimeout
	trans.conn = conn
	trans.BaseTransport = new(BaseTransport)
	trans.id = uuid.New()
	trans.connectTime = time.Now().Unix()
	trans.lastTime = 0
	return trans
}

func (t *TcpTransport) String() string {
	return "id:" + t.id + " addr:" + t.conn.RemoteAddr().String()
}

//实现 io.Writer
func (t *TcpTransport) Write(b []byte) (int, error) {
	t.lastTime = time.Now().Unix()
	return t.conn.Write(b)
}

//实现 io.Reader
func (t *TcpTransport) Read(b []byte) (int, error) {
	cnt, err := t.conn.Read(b)
	// if err != nil {
	// 	if e, ok := err.(net.Error); ok && e.Timeout() { //验证错误类型是否是超时错误
	// 		// Write space ping
	// 		if _, err := t.Write([]byte(" ")); err != nil {
	// 			return 0, err
	// 		}
	// 		t.SetReadTimeout()
	// 		return 0, nil
	// 	}
	// } else {
	// 	if t.hasSetReadTimeout {
	// 		t.SetReadTimeout()
	// 	}
	// }
	return cnt, err
}

func (t *TcpTransport) RemoteAddr() net.Addr {
	return t.conn.RemoteAddr()
}

func (t *TcpTransport) Close() error {
	return t.conn.Close()
}

func (t *TcpTransport) Id() string {
	return t.id
}

func (t *TcpTransport) Conn() net.Conn {
	return t.conn
}

func (t *TcpTransport) SetReadTimeout() error {
	err := t.conn.SetReadDeadline(time.Now().Add(t.readTimeout))
	if err != nil {
		return err
	}
	t.hasSetReadTimeout = true
	return nil
}

func (t *TcpTransport) UnsetReadTimeout() error {
	err := t.conn.SetReadDeadline(time.Time{})
	if err != nil {
		return err
	}

	t.hasSetReadTimeout = false
	return nil
}

func (t *TcpTransport) SetWriteTimeout() error {
	err := t.conn.SetWriteDeadline(time.Now().Add(t.writeTimeout))
	if err != nil {
		return err
	}
	return nil
}

func (t *TcpTransport) UnsetWriteTimeout() error {
	err := t.conn.SetWriteDeadline(time.Time{})
	if err != nil {
		return err
	}
	return nil
}

func (t *TcpTransport) ConnectTime() int64 {
	return t.connectTime
}

func (t *TcpTransport) LastTime() int64 {
	return t.lastTime
}
