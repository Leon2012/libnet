package libnet

type TransportHandler interface {
	OnClose(trans Transport)
	OnReceive(trans Transport, data []byte)
	OnConnect(trans Transport)
}
