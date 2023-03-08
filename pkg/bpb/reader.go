package bpb

import (
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
func (b *BytePacketBuffer) ReadQName(qName *string) error {
	var (
		currPosition uint = b.Pos()
		hasJumped    bool
		maxJumps     int = 5
		jumpCount    int
		delimeter    string
		strBuffer    []uint8
	)

	for {
		if jumpCount > maxJumps {
			return fmt.Errorf("limit of %v jumps exceeded", maxJumps)
		}
		len, err := b.Get(currPosition)
		if err != nil {
			return err
		}
		// Due to limitations to the size to a DNS packet to 512, as suggested in RFC 1035 there is
		// a kind of compression mechanism implemented.
		// When data like "google.com" appear multiple times in a packet, after the first occurence
		// they are referenced using the jump directive, which tells the parser to jump to the location pointed
		// and start reading the data for QName.
		// This jump directive is such that two MSBs of length label is set in ourcase 0xC0
		if (len & JumpTrigger) == JumpTrigger {
			if !hasJumped {
				b.Seek(currPosition + 2)
			}
			jumpTo, err := b.Get(currPosition + 1)
			if err != nil {
				return err
			}
			jumpOffset := (uint16(len)^JumpTrigger)<<8 | uint16(jumpTo)
			currPosition = uint(jumpOffset)
			hasJumped = true
			jumpCount++
			continue
		}

		currPosition++
		if len == 0 {
			break
		}
		*qName = fmt.Sprintf("%s%s", *qName, delimeter)
		strBuffer, err = b.GetRange(currPosition, uint(len))
		if err != nil {
			return err
		}
		*qName = fmt.Sprintf("%s%s",
			*qName,
			strings.ToLower(string(strBuffer)),
		)
		delimeter = "."
		currPosition += uint(len)
	}

	if !hasJumped {
		b.Seek(currPosition)
	}

	return nil
}

// func (b *BytePacketBuffer) ReadQName(outstr *string) error {

// 	var err error
// 	var pos = b.Pos()

// 	jumped := false
// 	max_jumps := 5
// 	jumps_performed := 0

// 	delim := ""
// 	for {

// 		if jumps_performed > max_jumps {
// 			return fmt.Errorf("limit of %v jumps exceeded", max_jumps)
// 		}

// 		var len byte
// 		if len, err = b.Get(pos); err != nil {
// 			return err
// 		}
// 		if len == 0xC0 {

// 			if !jumped {
// 				b.Seek(pos + 2)
// 			}

// 			var p_b2 byte
// 			if p_b2, err = b.Get(pos + 1); err != nil {
// 				return err
// 			}
// 			b2 := uint16(p_b2)

// 			offset := ((uint16(len) ^ 0xC0) << 8) | b2
// 			pos = uint(offset)

// 			jumped = true
// 			jumps_performed += 1
// 			continue
// 		} else {
// 			pos = pos + 1

// 			if len == 0 {
// 				break
// 			}

// 			*outstr = fmt.Sprintf("%s%s", *outstr, delim)

// 			var str_buffer []uint8

// 			if str_buffer, err = b.GetRange(pos, uint(len)); err != nil {
// 				return err
// 			}

// 			*outstr = fmt.Sprintf("%s%s", *outstr, strings.ToLower(string(str_buffer)))

// 			delim = "."

// 			pos = pos + uint(len)
// 		}
// 	}

// 	if !jumped {
// 		b.Seek(pos)

// 	}

// 	return nil
// }
