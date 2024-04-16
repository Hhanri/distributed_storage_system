package main

import (
	"log"

	"github.com/Hhanri/distributed_storage_system/p2p"
)

func main() {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddress: ":3000",
		Handshaker:    p2p.NOPHandshaker{},
		Decoder:       p2p.DefaultDecoder{},
	}

	transport := p2p.NewTCPTransport(tcpOpts)

	if err := transport.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
}
