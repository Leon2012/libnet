package libnet

import (
	"fmt"
	"io"
	"testing"
)

type MyClientCodec struct {
}

func (m *MyClientCodec) Encode(message []byte) ([]byte, error) {
	return message, nil
}

func (m *MyClientCodec) Decode(r io.Reader) ([]byte, error) {
	buf := make([]byte, 1024)
	if _, err := r.Read(buf); nil != err {
		return nil, err
	}
	return buf, nil
}

func TestClient(t *testing.T) {
	addr := "127.0.0.1:3333"
	client := NewTcpClient(addr, &MyClientCodec{}, 0, 0)
	err := client.Connect()
	if err != nil {
		fmt.Println(err)
	} else {
		n, _ := client.Send([]byte("hello"))
		fmt.Printf("send data length : %d \n", n)

		data, _ := client.Recv()
		fmt.Println("recv:" + string(data))
		client.Close()
	}
}
