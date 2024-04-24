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

	peekBuff := make([]byte, 1)

	if _, err := reader.Read(peekBuff); err != nil {
		return nil
	}

	// in case of a stream,
	// we do not encode what is sent over the network
	// we just set rpc.Stream to true to handle it in logic
	stream := peekBuff[0] == IncomingStream
	if stream {
		rpc.Stream = true
		return nil
	}

	buff := make([]byte, 1024)
	n, err := reader.Read(buff)

	if err != nil {
		return err
	}

	rpc.Payload = buff[:n]

	return nil
}
