package libnet

import (
	"io"
)

type Codec interface {
	Encode(message []byte) ([]byte, error)
	Decode(r io.Reader) ([]byte, error)
}
