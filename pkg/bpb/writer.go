package bpb

import (
	"fmt"
	"strings"
)

func (b *BytePacketBuffer) Write(val uint8) error {
	if b.pos >= b.size {
		return EOBError
	}

	b.Buf[b.pos] = val
	b.pos += 1
	return nil

}

func (b *BytePacketBuffer) WriteOneByte(val uint8) error {
	if err := b.Write(val); err != nil {
		return err
	}
	return nil
}

func (b *BytePacketBuffer) WriteTwoBytes(val uint16) error {

	if err := b.Write(uint8(val >> 8)); err != nil {
		return err
	}
	if err := b.Write(uint8(val & 0xFF)); err != nil {
		return err
	}
	return nil
}

func (b *BytePacketBuffer) WriteFourBytes(val uint32) error {
	if err := b.Write(uint8((val >> 24) & 0xFF)); err != nil {
		return err
	}
	if err := b.Write(uint8((val >> 16) & 0xFF)); err != nil {
		return err
	}
	if err := b.Write(uint8((val >> 8) & 0xFF)); err != nil {
		return err
	}
	if err := b.Write(uint8((val >> 0) & 0xFF)); err != nil {
		return err
	}
	return nil
}

func (b *BytePacketBuffer) WriteQName(q_name *string) error {
	for _, label := range strings.Split(*q_name, ".") {
		len := len(label)

		if len > 0x3f {
			return fmt.Errorf("single label is %v, which exceeds 63 characters of length", len)
		}

		if err := b.Write(uint8(len)); err != nil {
			return err
		}

		for _, x := range []byte(label) {
			if err := b.WriteOneByte(x); err != nil {
				return err
			}
		}

	}

	if err := b.WriteOneByte(0); err != nil {
		return err
	}

	return nil
}
