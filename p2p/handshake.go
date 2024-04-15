package p2p

import "errors"

var ErrorInvalidHandShake = errors.New("invalid hanshake")

type Handshaker interface {
	ShakeHands(peer Peer) error
}

type NOPHandshaker struct{}

func (hs NOPHandshaker) ShakeHands(peer Peer) error {
	return nil
}
