package header

import (
	"errors"
	"fmt"
	"strings"
)

type BytePacketBuffer struct {
	Buf [512]byte
	pos uint
}

func NewBuffer() BytePacketBuffer {
	b := BytePacketBuffer{
		Buf: [512]uint8{0},
		pos: 0,
	}
	return b
}

func (b *BytePacketBuffer) Pos() *uint {
	return &b.pos
}

func (b *BytePacketBuffer) Step(steps uint) error {
	b.pos += steps
	return nil
}

func (b *BytePacketBuffer) Seek(pos uint) error {
	b.pos = pos
	return nil
}

func (b *BytePacketBuffer) Read() (*uint8, error) {
	if b.pos >= 512 {
		return nil, EOBError
	}
	cur_byte := b.Buf[b.pos]
	// fmt.Println("before read", b.pos)
	b.pos += 1
	// fmt.Println("after read", b.pos)

	return &cur_byte, nil
}

func (b *BytePacketBuffer) Get(pos uint) (*uint8, error) {
	if b.pos >= 512 {
		return nil, EOBError
	}
	cur_buf := b.Buf[b.pos]

	return &cur_buf, nil
}

func (b *BytePacketBuffer) Get_range(start uint, len uint) ([]uint8, error) {
	if start+len >= 512 {
		return nil, EOBError
	}
	return b.Buf[start : start+len], nil
}

func (b *BytePacketBuffer) Read_u16() (*uint16, error) {

	byte1, err1 := b.Read()
	byte2, err2 := b.Read()
	if err1 != nil || err2 != nil {
		return nil, EOBError
	}

	// fmt.Println(*byte1)
	// fmt.Println(*byte2)

	two_bytes := (uint16(*byte1) << 8) | (uint16(*byte2))
	// fmt.Println(two_bytes)
	return &two_bytes, nil
}

func (b *BytePacketBuffer) Read_u32() (*uint32, error) {

	byte1, err1 := b.Read()
	byte2, err2 := b.Read()
	byte3, err3 := b.Read()
	byte4, err4 := b.Read()
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return nil, EOBError
	}

	four_bytes := (uint32(*byte1) << 24) | (uint32(*byte2) << 16) | (uint32(*byte3) << 8) | (uint32(*byte4))
	return &four_bytes, nil
}

func (b *BytePacketBuffer) Read_qname(outstr *string) error {

	var err error
	pos := b.Pos()
	var set_pos = *pos
	// b.Read()

	jumped := false
	max_jumps := 5
	jumps_performed := 0

	delim := ""
	for true {

		if jumps_performed > max_jumps {
			return errors.New(fmt.Sprintf("Limit of %v jumps exceeded", max_jumps))
		}

		var p_len *byte
		if p_len, err = b.Get(*pos); err != nil {
			return err
		}
		len := *p_len
		fmt.Println("len--", len)
		fmt.Println("test--", (len&0xC0) == 0xC0)

		if (len & 0xC0) == 0xC0 {

			if !jumped {
				// newV = *pos
				b.Seek(*pos + 2)
			}

			var p_b2 *byte
			if p_b2, err = b.Get(set_pos + 1); err != nil {
				return err
			}
			b2 := uint16(*p_b2)
			// fmt.Println("ab2", b2)

			offset := ((uint16(len) ^ 0xC0) << 8) | b2
			*pos = uint(offset)
			fmt.Println("offset----", set_pos)

			jumped = true
			jumps_performed += 1
			continue
		} else {
			// fmt.Println("position b", b.pos)
			*pos = *pos + 1
			// fmt.Println("position a", b.pos)

			if len == 0 {
				break
			}

			// fmt.Println("<<<<<<<>>>>>>>")
			// fmt.Println("len", uint(len))
			// fmt.Println("position", b.pos)
			*outstr = fmt.Sprintf("%s%s", *outstr, delim)

			var str_buffer []uint8

			if str_buffer, err = b.Get_range(*pos, uint(len)); err != nil {
				return err
			}

			// fmt.Println("Outstr> ", strings.ToLower(string(str_buffer)))

			*outstr = fmt.Sprintf("%s%s", *outstr, strings.ToLower(string(str_buffer)))

			delim = "."

			*pos = *pos + uint(len)
		}
	}
	// fmt.Println("out", *outstr)

	if !jumped {
		b.Seek(*pos)

	}

	return nil
}
