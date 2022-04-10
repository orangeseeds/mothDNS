package utils

import (
	"log"
	"net"

	"github.com/orangeseeds/DNSserver/core"
)

func ConstrPacket(id uint16, isRec bool, nameQtypes map[string]core.QueryType) core.DnsPacket {
	packet := core.NewPacket()
	packet.Header.Id = id

	for name, qtype := range nameQtypes {
		packet.Questions = append(packet.Questions, core.NewQuestion(name, qtype))
	}

	packet.Header.Questions = uint16(len(nameQtypes))
	return packet
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

func Lookup(nameQtypes map[string]core.QueryType, id uint16, serverType string, host string, port string) (*core.DnsPacket, error) {
	socket, err := net.Dial(serverType, host+":"+port)
	if err != nil {
		return nil, err
	}
	defer socket.Close()
	log.Printf("Looking up Ip addresses from %v", host)

	packet := ConstrPacket(id, true, nameQtypes)
	buffer := PacketToBuf(packet)

	_, err = socket.Write(From_512buffer(buffer.Buf))
	if err != nil {
		return nil, err
	}

	replyBuffer := make([]byte, 512)
	_, err = socket.Read(replyBuffer)
	if err != nil {
		return nil, err
	}

	packetBuffer := core.NewBuffer()
	packetBuffer.Buf = To_512buffer(replyBuffer)
	replyPacket, err := BufToPacket(packetBuffer)
	if err != nil {
		return nil, err
	}

	return replyPacket, nil
}

func CheckQuestions(questions []core.DnsQuestion) bool {
	for _, question := range questions {
		if question.Name == "127.0.0.1" {
			log.Printf("Question asking for %v, is not a valid question.", question.Name)
			return false
		}
	}
	return true
}
