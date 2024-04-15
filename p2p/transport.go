package p2p

// Peer is any remote note
type Peer interface {
	Close() error
}

// Transport is anything that handles communication
// between the nodes in a network (TCP, UDP, websockets, ...)
type Transport interface {
	ListenAndAccept() error
}
