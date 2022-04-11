package utils

import (
	"fmt"
	"log"
	"net"

	"github.com/orangeseeds/DNSserver/core"
)

func ConstrPacket(id uint16, isRec bool, nameQtypes map[string]core.QueryType) core.DnsPacket {
	packet := core.NewPacket()
	packet.Header.Id = id
	packet.Header.Recursion_desired = isRec

	for name, qtype := range nameQtypes {
		packet.Questions = append(packet.Questions, core.NewQuestion(name, qtype))
	}
	packet.Header.Questions = uint16(len(nameQtypes))
	return packet
}

func Lookup(nameQtypes map[string]core.QueryType, id uint16, serverType string, host string, port string) (*core.DnsPacket, error) {
	socket, err := net.Dial(serverType, host+":"+port)
	if err != nil {
		return nil, err
	}
	defer socket.Close()
	packet := ConstrPacket(id, true, nameQtypes)
	buffer := PacketToBuf(packet)

	_, err = socket.Write(From_512buffer(buffer.Buf))
	if err != nil {
		return nil, fmt.Errorf("error while writing to %v, %v", host, err)
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
