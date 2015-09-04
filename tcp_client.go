package libnet

import (
	"net"
)

type TcpClient struct {
	addr         string
	trans        Transport
	transHandler TransportHandler
	codec        Codec
}

func NewTcpClient(addr string, th TransportHandler, c Codec) *TcpClient {
	return &TcpClient{
		addr:         addr,
		transHandler: th,
		codec:        c,
		trans:        nil,
	}
}

func (t *TcpClient) Connect() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", t.addr)
	if err != nil {
		return err
	}

	conn, err := net.Dial("tcp", tcpAddr.String())
	if err != nil {
		return err
	}
}
