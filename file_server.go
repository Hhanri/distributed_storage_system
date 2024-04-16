package main

import (
	"fmt"

	"github.com/Hhanri/distributed_storage_system/p2p"
	"github.com/Hhanri/distributed_storage_system/store"
)

type FileServerOpts struct {
	store.StoreOpts
	Transport p2p.Transport
}

type FileServer struct {
	FileServerOpts
	store *store.Store

	quitCh chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	return &FileServer{
		FileServerOpts: opts,
		store:          store.NewStore(opts.StoreOpts),
		quitCh:         make(chan struct{}),
	}
}

func (fs *FileServer) loop() {

	defer func() {
		fmt.Println("File server shutting down")
		fs.Transport.Close()
	}()

	for {
		select {
		case msg := <-fs.Transport.Consume():
			fmt.Println(msg)
		case <-fs.quitCh:
			return
		}
	}

}

func (fs *FileServer) Stop() {
	close(fs.quitCh)
}

func (fs *FileServer) Start() error {
	if err := fs.Transport.ListenAndAccept(); err != nil {
		return err
	}

	fs.loop()

	return nil
}
