package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	tcpOpts := TCPTransportOpts{
		ListenAddress: ":3000",
		Handshaker:    NOPHandshaker{},
		Decoder:       GOBDecoder{},
	}
	transport := NewTCPTransport(tcpOpts)

	assert.Equal(
		t,
		transport,
		&TCPTransport{TCPTransportOpts: tcpOpts},
	)
	assert.Nil(t, transport.ListenAndAccept())
}
