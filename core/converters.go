package core

// Some utility function to reuse in the code
func PacketToBuf(p DnsPacket) BytePacketBuffer {
	buffer := NewBuffer()
	p.Write(&buffer)
	return buffer
}

func BufToPacket(b BytePacketBuffer) (*DnsPacket, error) {
	packet := NewPacket()
	_, err := packet.From_buffer(&b)
	if err != nil {
		return nil, err
	}
	return &packet, nil
}

func ConstrPacket(id uint16, isRec bool, nameQtypes map[string]QueryType) DnsPacket {
	packet := NewPacket()
	packet.Header.Id = id
	packet.Header.Recursion_desired = isRec

	for name, qtype := range nameQtypes {
		packet.Questions = append(packet.Questions, NewQuestion(name, qtype))
	}
	packet.Header.Questions = uint16(len(nameQtypes))
	return packet
}
