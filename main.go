package main

import (
	"bytes"
	"fmt"
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
			Root:          "./storage_content/" + listenAddr + "_network",
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

	time.Sleep(time.Second * 2)

	go func() {
		log.Fatal(server2.Start())
	}()

	time.Sleep(time.Second * 2)

	for i := 0; i < 10; i++ {
		data := bytes.NewReader([]byte("My big data file here!"))
		server2.StoreData(fmt.Sprintf("myprivatekey_%d", i), data)
		time.Sleep(time.Millisecond * 5)
	}

	// reader, err := server2.GetData("myprivatekey")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// bytes, err := io.ReadAll(reader)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(string(bytes))
	select {}
}
