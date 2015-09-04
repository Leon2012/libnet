package libnet

import (
	"net"
	"time"
)

type TcpClient struct {
	addr               string
	trans              Transport
	codec              Codec
	rTimeout, wTimeout time.Duration
}

func NewTcpClient(addr string, c Codec, rTimeout, wTimeout int) *TcpClient {
	return &TcpClient{
		addr:     addr,
		codec:    c,
		trans:    nil,
		rTimeout: time.Duration(rTimeout) * time.Second,
		wTimeout: time.Duration(wTimeout) * time.Second,
	}
}

func (t *TcpClient) Connect() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", t.addr)
	if err != nil {
		return err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}

	t.trans = NewTCPTransport(conn, t.rTimeout, t.wTimeout)
	t.trans.Conn().SetDeadline(time.Now().Add(1e9))
	return nil
}

func (t *TcpClient) Send(message []byte) (int, error) {
	buf, err := t.codec.Encode(message)
	if err != nil {
		return 0, err
	} else {
		return t.trans.Write(buf)
	}
}

func (t *TcpClient) Recv() ([]byte, error) {
	buf, err := t.codec.Decode(t.trans.Conn())
	if err != nil {
		return nil, err
	} else {
		return buf, nil
	}
}

func (t *TcpClient) Close() {
	t.trans.Conn().SetDeadline(time.Now())
	t.trans.Close()
}
