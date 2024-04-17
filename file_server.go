package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/Hhanri/distributed_storage_system/p2p"
	"github.com/Hhanri/distributed_storage_system/store"
)

type FileServerOpts struct {
	store.StoreOpts
	Transport      p2p.Transport
	BootstrapNodes []string
}

type FileServer struct {
	FileServerOpts

	peerslock sync.Mutex
	peers     map[string]p2p.Peer

	store  *store.Store
	quitCh chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	return &FileServer{
		FileServerOpts: opts,
		store:          store.NewStore(opts.StoreOpts),
		quitCh:         make(chan struct{}),

		peerslock: sync.Mutex{},
		peers:     make(map[string]p2p.Peer),
	}
}

func (fs *FileServer) OnPeer(peer p2p.Peer) error {
	fs.peerslock.Lock()
	defer fs.peerslock.Unlock()

	fs.peers[peer.RemoteAddr().String()] = peer

	log.Printf("connected with remote %s", peer.RemoteAddr())

	return nil
}

func (fs *FileServer) loop() {

	defer func() {
		fmt.Println("File server shutting down")
		fs.Transport.Close()
	}()

	for {
		select {
		case msg := <-fs.Transport.Consume():
			m := &Message{}
			if err := gob.NewDecoder(bytes.NewReader(msg.Payload)).Decode(m); err != nil {
				log.Println(err)
			}

			if err := fs.handleMessage(m); err != nil {
				log.Println(err)
			}
		case <-fs.quitCh:
			return
		}
	}

}

func (fs *FileServer) handleMessage(msg *Message) error {
	fmt.Printf("Received data: %+v\n", msg.Payload)
	return nil
}

func (fs *FileServer) broadcast(p *Message) error {
	peers := []io.Writer{}
	for _, peer := range fs.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...)

	return gob.NewEncoder(mw).Encode(*p)
}

func (fs *FileServer) StoreData(key string, reader io.Reader) error {
	buff := new(bytes.Buffer)
	tee := io.TeeReader(reader, buff)

	if err := fs.store.Write(key, tee); err != nil {
		return err
	}

	payload := &MessageData{
		Key:  key,
		Data: buff.Bytes(),
	}

	return fs.broadcast(&Message{
		From:    "TODO",
		Payload: *payload,
	})
}

func (fs *FileServer) Stop() {
	close(fs.quitCh)
}

func (fs *FileServer) Start() error {
	if err := fs.Transport.ListenAndAccept(); err != nil {
		return err
	}

	fs.bootstrapNetwork()
	fs.loop()

	return nil
}

func (fs *FileServer) bootstrapNetwork() {
	for _, addr := range fs.BootstrapNodes {
		go func(addr string) {
			if err := fs.Transport.Dial(addr); err != nil {
				log.Println("dial error: ", err)
			}
		}(addr)
	}
}
