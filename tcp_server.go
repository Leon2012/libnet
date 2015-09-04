package libnet

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

const (
	DEFAULT_MAX_BUFFER_LENGTH          = 1024
	DEFAULT_MAX_PACKAGE_CHANNEL_LENGTH = 1000
)

type TcpServer struct {
	addr               string
	ch                 chan bool
	bufLen             int
	packChanLen        int
	waitGroup          *sync.WaitGroup
	connections        map[string]Transport
	listener           *net.TCPListener
	rTimeout, wTimeout time.Duration
	transHandler       TransportHandler
	serverHandler      TcpServerHandler
	codec              Codec
}

func NewServer(ip string, port int, th TransportHandler, sh TcpServerHandler, c Codec) *TcpServer {
	addr := fmt.Sprintf("%s:%d", port)
	return NewServerWithAddr(addr, th, sh, c)
}

func NewServerWithAddr(addr string, th TransportHandler, sh TcpServerHandler, c Codec) *TcpServer {
	return NewServerWithTimeout(addr, th, sh, c, 0, 0)
}

func NewServerWithTimeout(addr string, th TransportHandler, sh TcpServerHandler, c Codec, wTimeout, rTimeout int) *TcpServer {
	serv := &TcpServer{
		ch:            make(chan bool),
		waitGroup:     &sync.WaitGroup{},
		connections:   make(map[string]Transport),
		listener:      nil,
		addr:          addr,
		rTimeout:      time.Duration(rTimeout) * time.Second,
		wTimeout:      time.Duration(wTimeout) * time.Second,
		bufLen:        DEFAULT_MAX_BUFFER_LENGTH,
		packChanLen:   DEFAULT_MAX_PACKAGE_CHANNEL_LENGTH,
		transHandler:  th,
		serverHandler: sh,
		codec:         c,
	}
	return serv
}

func (s *TcpServer) Start() {
	lis, err := s.listen()
	if err != nil {
		log.Println("listen error : " + err.Error())
		//return err
	} else {
		s.listener = lis
		s.waitGroup.Add(1)
		go s.accept(lis)
		if s.serverHandler != nil {
			s.serverHandler.OnStart(s)
		}
	}
}

func (s *TcpServer) Stop() {
	if s.listener != nil {
		s.listener.SetDeadline(time.Now()) //close tcp accept
		for _, trans := range s.connections {
			//log.Println("id   : " + id)
			trans.Close()
		}
	}
	s.ch <- true
	//close(s.ch)
	s.waitGroup.Wait()
	if s.serverHandler != nil {
		s.serverHandler.OnStop(s)
	}
}

func (s *TcpServer) Send(id string, message []byte) (int, error) {
	trans, ok := s.connections[id]
	if !ok {
		return 0, errors.New("cannot find id : " + id)
	} else {
		buf, err := s.codec.Encode(message)
		if err != nil {
			return 0, err
		} else {
			return trans.Write(buf)
		}
	}
}

func (s *TcpServer) Exist(id string) bool {
	_, ok := s.connections[id]
	if !ok {
		return false
	} else {
		return true
	}
}

func (s *TcpServer) SetBufferLength(length int) {
	s.bufLen = length
}

func (s *TcpServer) SetPackageChannelLength(length int) {
	s.packChanLen = length
}

func (s *TcpServer) GetTransport(id string) Transport {
	trans, ok := s.connections[id]
	if !ok {
		return nil
	}
	return trans
}

func (s *TcpServer) GetAddr() string {
	return s.addr
}

func (s *TcpServer) listen() (*net.TCPListener, error) {
	//s.addr = addr
	var err error
	var lis *net.TCPListener
	lis, err = listenTCP(s.addr)
	if err != nil {
		log.Println("listen tcp error:", err)
		return nil, err
	}
	//s.listener = lis
	return lis, nil
}

func (s *TcpServer) accept(lis *net.TCPListener) {
	defer lis.Close()
	defer s.waitGroup.Done()
	for {
		select {
		case <-s.ch: //stop goroutine
			//log.Println("close tcp listener")
			return
		default:
			//log.Println("close 11111 .....")
		}
		var trans Transport = nil
		conn, err := lis.AcceptTCP()
		if nil != err {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				//log.Println("Stop accepting connections")
				continue
			}
			log.Println(err)
		}
		conn.SetLinger(-1)
		trans = NewTCPTransport(conn, s.rTimeout, s.wTimeout)
		s.connections[trans.Id()] = trans
		if s.transHandler != nil {
			s.transHandler.OnConnect(trans)
		}
		s.waitGroup.Add(1)
		go s.run(trans)
	}
}

func (s *TcpServer) run(trans Transport) {
	defer trans.Close()
	defer s.waitGroup.Done()
	//addr := trans.RemoteAddr().String()
	receivePackets := make(chan []byte, s.packChanLen)
	chStop := make(chan bool) // 通知停止消息处理
	defer func() {            //出错处理
		defer func() {
			if e := recover(); e != nil {
				log.Printf("Panic: %v", e)
			}
		}()
		//log.Printf("Disconnect: %v", addr)
		chStop <- true
	}()
	// 处理接收到的包
	go s.handlePackets(trans, receivePackets, chStop)
	for {
		select {
		case <-s.ch:
			//log.Println("disconnecting ", trans.RemoteAddr().String())
			return
		default:
		}
		trans.Conn().SetReadDeadline(time.Now().Add(1e9))
		buf, err := s.codec.Decode(trans.Conn())
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() { //handle read timeout error
				continue
			}
			log.Println(err)
			return
		}
		receivePackets <- buf
	}
}

func (s *TcpServer) handlePackets(trans Transport, receivePackets <-chan []byte, chStop <-chan bool) {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("Panic: %v", e)
		}
	}()
	for {
		select {
		case <-chStop:
			trans.Close()
			//log.Printf("Stop handle receivePackets.")
			delete(s.connections, trans.Id())
			if s.transHandler != nil {
				s.transHandler.OnClose(trans)
			}
			return

		// 消息包处理
		case bytes := <-receivePackets:
			//log.Println("receive bytes : ", string(bytes))
			if s.transHandler != nil {
				s.transHandler.OnReceive(trans, bytes)
			}
		}
	}
}
