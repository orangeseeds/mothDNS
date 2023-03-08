package bpb

import (
	"errors"
	"fmt"
	"strings"
)

// Read returns the byte at the current position in the buffer and increases Pos by one.
func (b *BytePacketBuffer) Read() (uint8, error) {
	if b.pos >= b.size {
		return 0, EOBError
	}
	currByte := b.Buf[b.pos]
	b.pos += 1
	return currByte, nil
}

// Get returns the byte at the current position in the buffer without increasing the value of Pos.
func (b *BytePacketBuffer) Get(pos uint) (uint8, error) {
	if b.pos >= b.size {
		return 0, EOBError
	}
	return b.Buf[pos], nil
}

// GetRange returns a slice of Buf starting at start and of length len.
func (b *BytePacketBuffer) GetRange(start uint, len uint) ([]uint8, error) {
	if start+len >= b.size {
		return nil, EOBError
	}
	return b.Buf[start : start+len], nil
}

// ReadTwoBytes returns a uint16 value containing two bytes from current position and advances Pos by two.
func (b *BytePacketBuffer) ReadTwoBytes() (uint16, error) {

	byte1, err1 := b.Read()
	byte2, err2 := b.Read()
	if err1 != nil || err2 != nil {
		return 0, EOBError
	}
	twoBytes := (uint16(byte1) << 8) | (uint16(byte2))
	return twoBytes, nil
}

// ReadTwoBytes returns a uint32 value containing four bytes from current position and advances Pos by four.
func (b *BytePacketBuffer) ReadFourBytes() (uint32, error) {

	byte1, err1 := b.Read()
	byte2, err2 := b.Read()
	byte3, err3 := b.Read()
	byte4, err4 := b.Read()
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return 0, EOBError
	}
	fourBytes := (uint32(byte1) << 24) | (uint32(byte2) << 16) | (uint32(byte3) << 8) | (uint32(byte4))
	return fourBytes, nil
}

// TODO: refacor the ReadQName function
func (b *BytePacketBuffer) ReadName() (error, string, uint) {
	pos := b.Pos()
	len, err := b.Get(b.Pos())
	if err != nil {
		return errors.New("Idk some kinda error!"), "", 0
	}

	if (len & 0xC0) == 0xC0 {

		b.Seek(pos + 2)

		var p_b2 byte
		if p_b2, err = b.Get(pos + 1); err != nil {
			return err, "", 0
		}
		b2 := uint16(p_b2)

		offset := ((uint16(len) ^ 0xC0) << 8) | b2
		pos = uint(offset)
		return nil, "", pos
	}

	return nil, "", 0
}

func (b *BytePacketBuffer) ReadQName(outstr *string) error {

	var err error
	var pos = b.Pos()

	jumped := false
	max_jumps := 5
	jumps_performed := 0

	delim := ""
	for {

		if jumps_performed > max_jumps {
			return fmt.Errorf("limit of %v jumps exceeded", max_jumps)
		}

		var len byte
		if len, err = b.Get(pos); err != nil {
			return err
		}
		if (len & 0xC0) == 0xC0 {

			if !jumped {
				b.Seek(pos + 2)
			}

			var p_b2 byte
			if p_b2, err = b.Get(pos + 1); err != nil {
				return err
			}
			b2 := uint16(p_b2)

			offset := ((uint16(len) ^ 0xC0) << 8) | b2
			pos = uint(offset)

			jumped = true
			jumps_performed += 1
			continue
		} else {
			pos = pos + 1

			if len == 0 {
				break
			}

			*outstr = fmt.Sprintf("%s%s", *outstr, delim)

			var str_buffer []uint8

			if str_buffer, err = b.GetRange(pos, uint(len)); err != nil {
				return err
			}

			*outstr = fmt.Sprintf("%s%s", *outstr, strings.ToLower(string(str_buffer)))

			delim = "."

			pos = pos + uint(len)
		}
	}

	if !jumped {
		b.Seek(pos)

	}

	return nil
}
