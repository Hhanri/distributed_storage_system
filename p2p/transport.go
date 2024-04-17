package p2p

import "net"

// Peer is any remote note
type Peer interface {
	net.Conn
	Send([]byte) error
}

// Transport is anything that handles communication
// between the nodes in a network (TCP, UDP, websockets, ...)
type Transport interface {
	Dial(addr string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
