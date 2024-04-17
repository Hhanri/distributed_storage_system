package main

import (
	"bytes"
	"log"
	"time"

	"github.com/Hhanri/distributed_storage_system/p2p"
	"github.com/Hhanri/distributed_storage_system/store"
)

func makeServer(listenAddr string, nodes []string) *FileServer {
	tcpOpts := p2p.TCPTransportOpts{
		ListenAddress: listenAddr,
		Handshaker:    p2p.NOPHandshaker{},
		Decoder:       p2p.DefaultDecoder{},
	}
	transport := p2p.NewTCPTransport(tcpOpts)

	fileServerOtps := FileServerOpts{
		StoreOpts: store.StoreOpts{
			Root:          listenAddr + "_network",
			PathTransform: store.HashPathTransform,
		},
		BootstrapNodes: nodes,

		Transport: transport,
	}

	server := NewFileServer(fileServerOtps)
	transport.OnPeer = server.OnPeer

	return server
}

func main() {

	server1 := makeServer(":3000", []string{})
	server2 := makeServer(":4000", []string{":3000"})

	go func() {
		log.Fatal(server1.Start())
	}()

	time.Sleep(time.Second * 3)

	go func() {
		log.Fatal(server2.Start())
	}()

	time.Sleep(time.Second * 3)

	data := bytes.NewReader([]byte("My big data file here!"))
	server2.StoreData("myprivatekey", data)

	select {}
}
