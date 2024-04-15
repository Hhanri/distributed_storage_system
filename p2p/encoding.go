package p2p

import (
	"encoding/gob"
	"io"
)

type Decoder interface {
	Decode(reader io.Reader, msg *Message) error
}

type GOBDecoder struct{}

func (d GOBDecoder) Decode(reader io.Reader, msg *Message) error {
	return gob.NewDecoder(reader).Decode(msg)
}

type DefaultDecoder struct{}

func (d DefaultDecoder) Decode(reader io.Reader, msg *Message) error {
	buff := make([]byte, 1024)
	n, err := reader.Read(buff)

	if err != nil {
		return err
	}

	msg.Payload = buff[:n]

	return nil
}
