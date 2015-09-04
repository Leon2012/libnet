package libnet

import (
	_ "time"
)

type BaseTransport struct {
	id          string
	connectTime int64
	lastTime    int64
}
