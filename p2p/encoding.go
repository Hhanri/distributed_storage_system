package p2p

import (
	"encoding/gob"
	"io"
)

type Decoder interface {
	Decode(reader io.Reader, v any) error
}

type GOBDecoder struct{}

func (d GOBDecoder) Decode(reader io.Reader, v any) error {
	return gob.NewDecoder(reader).Decode(v)
}
