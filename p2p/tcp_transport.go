package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCPPeer represents the remote node over a TCP established connection
type TCPPeer struct {
	// underlying connection of the peer
	conn net.Conn

	// if server dials and retrieves a conn => outbound == true
	// if server accepts and retrieves a conn => outbound == false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

type TCPTransportOpts struct {
	ListenAddress string
	Handshaker    Handshaker
	Decoder       Decoder
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener

	rpcCh chan RPC
	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcCh:            make(chan RPC),
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	listener, err := net.Listen("tcp", t.ListenAddress)
	if err != nil {
		return err
	}

	t.listener = listener
	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcCh
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}

		fmt.Printf("New incoming connection: %v\n", conn)
		go t.handleConnection(conn)
	}
}

func (t *TCPTransport) handleConnection(conn net.Conn) {

	peer := NewTCPPeer(conn, true)

	if err := t.Handshaker.ShakeHands(peer); err != nil {
		fmt.Printf("Invalid handshake: %s\n", err)
		peer.Close()
		return
	}

	rpc := &RPC{}
	for {
		if err := t.Decoder.Decode(conn, rpc); err != nil {
			fmt.Printf("TCP error: %s\n", err)
			continue
		}
		rpc.From = conn.RemoteAddr()
		fmt.Printf("message: %+v\n", rpc)
	}

}
