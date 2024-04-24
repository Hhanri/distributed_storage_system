package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

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
		case rpc := <-fs.Transport.Consume():
			msg := &Message{}
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(msg); err != nil {
				log.Println(err)
				return
			}

			if err := fs.handleMessage(rpc.From.String(), msg); err != nil {
				log.Println(err)
			}

		case <-fs.quitCh:
			return
		}
	}

}

func (fs *FileServer) handleMessage(from string, msg *Message) error {
	switch v := msg.Payload.(type) {
	case MessageStoreFile:
		return fs.handleMessageStoreFile(from, &v)
	case MessageGetFile:
		return fs.handleMessageGetFile(from, &v)
	}

	return nil
}

func (fs *FileServer) handleMessageGetFile(from string, msg *MessageGetFile) error {
	if !fs.store.Has(msg.Key) {
		return fmt.Errorf("[%s] Needs to serve file (%s) but was not found on disk\n", fs.Transport.Addr(), msg.Key)
	}

	fmt.Printf("[%s] Serving file (%s) over the network\n", fs.Transport.Addr(), msg.Key)

	fileSize, r, err := fs.store.Read(msg.Key)
	if err != nil {
		return err
	}

	if rc, ok := r.(io.ReadCloser); ok {
		defer rc.Close()
	}

	peer, ok := fs.peers[from]
	if !ok {
		return fmt.Errorf("peer %s not found\n", from)
	}

	// First we send the "incomingStream" byte to the peer
	// then we can send the file size as a int64

	peer.Send([]byte{p2p.IncomingStream})
	binary.Write(peer, binary.LittleEndian, fileSize)
	n, err := io.Copy(peer, r)
	if err != nil {
		return err
	}

	fmt.Printf("[%s] Written %d bytes over the network to %s\n", fs.Transport.Addr(), n, from)

	return nil
}

func (fs *FileServer) handleMessageStoreFile(from string, msg *MessageStoreFile) error {
	peer, ok := fs.peers[from]
	defer peer.CloseStream()

	if !ok {
		return fmt.Errorf("peer (%s) could no be found", from)
	}
	fmt.Printf("[%s]:\n", fs.Transport.Addr())
	_, err := fs.store.Write(msg.Key, io.LimitReader(peer, msg.Size))
	if err != nil {
		return err
	}

	return nil
}

func (fs *FileServer) stream(p *Message) error {
	peers := []io.Writer{}
	for _, peer := range fs.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...)

	return gob.NewEncoder(mw).Encode(*p)
}

func (fs *FileServer) broadcast(msg *Message) error {
	msgBuff := new(bytes.Buffer)
	if err := gob.NewEncoder(msgBuff).Encode(msg); err != nil {
		return err
	}

	for _, peer := range fs.peers {
		peer.Send([]byte{p2p.IncomingMessage})
		if err := peer.Send(msgBuff.Bytes()); err != nil {
			return err
		}
	}
	return nil
}

func (fs *FileServer) StoreData(key string, reader io.Reader) error {

	fileBuff := new(bytes.Buffer)
	tee := io.TeeReader(reader, fileBuff)

	size, err := fs.store.Write(key, tee)
	if err != nil {
		return err
	}

	msg := Message{
		Payload: MessageStoreFile{
			Key:  key,
			Size: size,
		},
	}

	if err := fs.broadcast(&msg); err != nil {
		return err
	}

	time.Sleep(time.Millisecond * 5)

	for _, peer := range fs.peers {
		peer.Send([]byte{p2p.IncomingStream})
		_, err := io.Copy(peer, fileBuff)
		if err != nil {
			return err
		}
	}

	return nil
}

func (fs *FileServer) GetData(key string) (io.Reader, error) {
	if fs.store.Has(key) {
		fmt.Printf("[%s] Serving file (%s) from local disk\n", fs.Transport.Addr(), key)
		_, r, err := fs.store.Read(key)
		return r, err
	}

	fmt.Printf("[%s] File (%s) not found locally, fetching from network\n", fs.Transport.Addr(), key)

	msg := Message{
		Payload: MessageGetFile{
			Key: key,
		},
	}

	if err := fs.broadcast(&msg); err != nil {
		return nil, err
	}

	time.Sleep(time.Millisecond * 5)

	for _, peer := range fs.peers {

		// First read the file size so we can limit the amount of bytes to read from the connection
		// so it will not keep hanging

		var fileSize int64
		binary.Read(peer, binary.LittleEndian, &fileSize)

		n, err := fs.store.Write(key, io.LimitReader(peer, fileSize))
		if err != nil {
			return nil, err
		}

		fmt.Printf("[%s] Received %d bytes over the nentwork from [%s]\n", fs.Transport.Addr(), n, peer.RemoteAddr())
		peer.CloseStream()
	}

	_, r, err := fs.store.Read(key)
	return r, err
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

func init() {
	gob.Register(MessageStoreFile{})
	gob.Register(MessageGetFile{})
}
