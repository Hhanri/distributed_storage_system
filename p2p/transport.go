package p2p

import "net"

// Peer is any remote note
type Peer interface {
	Send([]byte) error
	RemoteAddr() net.Addr
	Close() error
}

// Transport is anything that handles communication
// between the nodes in a network (TCP, UDP, websockets, ...)
type Transport interface {
	Dial(addr string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
