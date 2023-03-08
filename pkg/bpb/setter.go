package bpb

func (b *BytePacketBuffer) Set(pos uint, val uint8) error {
	if pos >= b.size {
		return EOBError
	}
	b.Buf[pos] = val
	return nil
}

func (b *BytePacketBuffer) SetTwoBytes(pos uint, val uint16) error {
	if err := b.Set(pos, uint8(val>>8)); err != nil {
		return err
	}
	if err := b.Set(pos+1, uint8(val&0xFF)); err != nil {
		return err
	}
	return nil
}
