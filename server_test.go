package libnet

import (
	"fmt"
	//"log"
	"testing"
	//"time"
	"io"
)

type TransHandler struct {
}

func NewTransHandler() *TransHandler {
	return &TransHandler{}
}

func (t *TransHandler) OnClose(trans Transport) {
	fmt.Println(trans.Id() + " closed")
}

func (t *TransHandler) OnReceive(trans Transport, data []byte) {
	fmt.Println(data)
}

func (t *TransHandler) OnConnect(trans Transport) {
	fmt.Println(trans.Id() + " connected")
}

type MyTcpServerHandler struct {
}

func (t *MyTcpServerHandler) OnStart(s *TcpServer) {
	fmt.Println(s.GetAddr() + " onStart")
}

func (t *MyTcpServerHandler) OnStop(s *TcpServer) {
	fmt.Println(s.GetAddr() + " onStop")
}

type MyCodec struct {
}

func (m *MyCodec) Encode(message []byte) ([]byte, error) {
	return message, nil
}

func (m *MyCodec) Decode(r io.Reader) ([]byte, error) {
	buf := make([]byte, 8)
	if _, err := r.Read(buf); nil != err {
		return nil, err
	}
	return buf, nil
}

func TestServer(t *testing.T) {
	addr := "127.0.0.1:5556"
	th := NewTransHandler()
	serv := NewServerWithAddr(addr, th, &MyTcpServerHandler{}, &MyCodec{})
	serv.SetBufferLength(8)
	serv.Start()
	fmt.Println("start server......")
	SignalHandle()
	serv.Stop()

}
