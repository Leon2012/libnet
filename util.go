package libnet

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func listenTCP(addr string) (*net.TCPListener, error) {
	laddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	listener, err := net.ListenTCP("tcp", laddr)
	if nil != err {
		return nil, err
	}

	return listener, nil
}

func listenUDP(addr string) (*net.UDPConn, error) {
	laddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	listener, err := net.ListenUDP("udp", laddr)
	if nil != err {
		return nil, err
	}

	return listener, nil
}

func SignalHandle() {
	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal, 1)
	over := make(chan bool, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-ch
		log.Println(sig)
		over <- true
	}()

	log.Println(<-over)
}
