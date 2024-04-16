package main

import (
	"log"
	"time"

	"github.com/Hhanri/distributed_storage_system/p2p"
	"github.com/Hhanri/distributed_storage_system/store"
)

func main() {

	tcpOpts := p2p.TCPTransportOpts{
		ListenAddress: ":3000",
		Handshaker:    p2p.NOPHandshaker{},
		Decoder:       p2p.DefaultDecoder{},
	}
	transport := p2p.NewTCPTransport(tcpOpts)

	fileServerOtps := FileServerOpts{
		StoreOpts: store.StoreOpts{
			Root:          store.DefaultRootStorage,
			PathTransform: store.HashPathTransform,
		},
		BootstrapNodes: []string{":4000"},

		Transport: transport,
	}

	server := NewFileServer(fileServerOtps)

	go func() {
		time.Sleep(time.Second * 3)
		server.Stop()
	}()

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}

}
