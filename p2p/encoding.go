package p2p

import (
	"encoding/gob"
	"io"
)

type Decoder interface {
	Decode(reader io.Reader, rpc *RPC) error
}

type GOBDecoder struct{}

func (d GOBDecoder) Decode(reader io.Reader, rpc *RPC) error {
	return gob.NewDecoder(reader).Decode(rpc)
}

type DefaultDecoder struct{}

func (d DefaultDecoder) Decode(reader io.Reader, rpc *RPC) error {
	buff := make([]byte, 1024)
	n, err := reader.Read(buff)

	if err != nil {
		return err
	}

	rpc.Payload = buff[:n]

	return nil
}
