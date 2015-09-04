package libnet

type TcpServerHandler interface {
	OnStart(s *TcpServer)
	OnStop(s *TcpServer)
}
