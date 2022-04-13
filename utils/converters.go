package utils

import (
	"github.com/orangeseeds/DNSserver/core"
)

func To_512buffer(buf []byte) [512]byte {
	reply_buff := [512]byte{0}
	// reply_buff := make([]byte, 512)
	for i := range buf {
		reply_buff[i] = buf[i]
	}

	return reply_buff
}

func From_512buffer(buf [512]byte) []byte {
	reply_buff := make([]byte, 250)
	for i := range reply_buff {
		reply_buff[i] = buf[i]
	}

	return reply_buff
}

func PacketToBuf(p core.DnsPacket) core.BytePacketBuffer {
	buffer := core.NewBuffer()
	p.Write(&buffer)
	return buffer
}

func BufToPacket(b core.BytePacketBuffer) (*core.DnsPacket, error) {
	packet := core.NewPacket()
	_, err := packet.From_buffer(&b)
	if err != nil {
		return nil, err
	}
	return &packet, nil
}
