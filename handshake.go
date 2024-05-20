package main

import (
	"fmt"
	"io"
)

// A Handshake is a special message that a peer uses to identify itself
type Handshake struct {
	Pstr     string // BitTorrent protocol
	InfoHash [20]byte
	PeerID   [20]byte
}

// New creates a new handshake with the standard pstr
func NewHandshake(infoHash, peerID [20]byte) *Handshake {
	return &Handshake{
		Pstr:     "BitTorrent protocol",
		InfoHash: infoHash,
		PeerID:   peerID,
	}
}

// Serialize serializes the handshake to a buffer
func (h *Handshake) Serialize() []byte {
	// Create a buffer with enough space to hold the serialized handshake
	buf := make([]byte, len(h.Pstr)+49)
	// The first byte indicates the length of the protocol identifier
	buf[0] = byte(len(h.Pstr))
	// Initialize a variable to keep track of the current position in the buffer
	curr := 1
	// Copy the protocol identifier (Pstr) into the buffer
	curr += copy(buf[curr:], h.Pstr)
	// Copy 8 reserved bytes (all set to 0) into the buffer
	curr += copy(buf[curr:], make([]byte, 8)) // 8 reserved bytes
	// Copy the infohash into the buffer
	curr += copy(buf[curr:], h.InfoHash[:])
	// Copy the peer ID into the buffer
	curr += copy(buf[curr:], h.PeerID[:])
	// Return the serialized buffer
	return buf
}

// Read parses a handshake from a stream
func ReadHandshake(r io.Reader) (*Handshake, error) {
	lengthBuf := make([]byte, 1)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, err
	}
	pstrlen := int(lengthBuf[0])

	if pstrlen == 0 {
		err := fmt.Errorf("pstrlen cannot be 0")
		return nil, err
	}

	handshakeBuf := make([]byte, 48+pstrlen)
	_, err = io.ReadFull(r, handshakeBuf)
	if err != nil {
		return nil, err
	}

	var infoHash, peerID [20]byte

	copy(infoHash[:], handshakeBuf[pstrlen+8:pstrlen+8+20])
	copy(peerID[:], handshakeBuf[pstrlen+8+20:])

	h := Handshake{
		Pstr:     string(handshakeBuf[0:pstrlen]),
		InfoHash: infoHash,
		PeerID:   peerID,
	}

	return &h, nil
}
