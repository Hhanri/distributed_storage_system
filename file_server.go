package main

import (
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
}

func NewFileServer(opts FileServerOpts) *FileServer {

	return &FileServer{
		FileServerOpts: opts,
		store:          store.NewStore(opts.StoreOpts),
	}
}

func (fs *FileServer) Start() error {
	return fs.Transport.ListenAndAccept()
}
