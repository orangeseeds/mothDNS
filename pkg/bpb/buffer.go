package bpb

// bpf => byte packet buffer

import (
	"errors"
)

var EOBError = errors.New("Buffer greater than or equals 512")

const (
	MaxSize     = 512
	JumpTrigger = 0xC0
)

// This struct is responsible for all operations related to creating
// & manipulating byte buffer for out UPD packets.
type BytePacketBuffer struct {
	// Buf contains the packet buffer
	Buf []byte `json:"buffer"`
	// pos keeps track of current position in the buffer.
	// This is used to check the point in the buffer up to which data has been read or written.
	pos uint
	// size contains the maximum size of Buf
	size uint
}

func New() BytePacketBuffer {
	b := BytePacketBuffer{
		pos:  0,
		size: MaxSize,
	}
	b.Buf = make([]byte, b.size)
	return b
}

func (b *BytePacketBuffer) Pos() uint {
	return b.pos
}

func (b *BytePacketBuffer) Size() uint {
	return b.size
}

func (b *BytePacketBuffer) Step(steps uint) error {
	b.pos += steps
	return nil
}

func (b *BytePacketBuffer) Seek(pos uint) error {
	b.pos = pos
	return nil
}
