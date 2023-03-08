package core

import (
	"github.com/orangeseeds/mothDNS/pkg/bpb"
)

// Some utility function to reuse in the code
func PacketToBuf(p DNSPacket) bpb.BytePacketBuffer {
	buffer := bpb.New()
	p.Write(&buffer)
	buffer.Buf = buffer.Buf[0 : buffer.Pos()+1]
	return buffer
}

func BufToPacket(b bpb.BytePacketBuffer) (*DNSPacket, error) {
	packet := NewPacket()
	_, err := packet.From_buffer(&b)
	if err != nil {
		return nil, err
	}
	return &packet, nil
}

func ConstrPacket(id uint16, isRec bool, nameQtypes map[string]QueryType) DNSPacket {
	packet := NewPacket()
	packet.Header.Id = id
	packet.Header.RecursionDesired = isRec

	for name, qtype := range nameQtypes {
		packet.Questions = append(packet.Questions, NewQuestion(name, qtype))
	}
	packet.Header.Questions = uint16(len(nameQtypes))
	return packet
}
