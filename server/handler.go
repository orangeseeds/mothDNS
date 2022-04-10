package server

import (
	"DNSserver/core"
	"DNSserver/utils"
	"log"
	"net"
)

func HandleConnection(socket net.PacketConn, addr net.Addr, buf []byte) {
	rcvBuffer := core.NewBuffer()
	rcvBuffer.Buf = utils.To_512buffer(buf)
	rcvPacket, err := utils.BufToPacket(rcvBuffer)
	if err != nil {
		log.Println("Invalid buffer from lookup response...")
	}
	if utils.CheckQuestions(rcvPacket.Questions) {
		socket.WriteTo(utils.From_512buffer(rcvBuffer.Buf), addr)
		return
	}
	var nameQtypes map[string]core.QueryType
	for _, question := range rcvPacket.Questions {
		nameQtypes[question.Name] = question.Qtype
	}

	respPacket, err := utils.Lookup(nameQtypes, rcvPacket.Header.Id, "udp", "8.8.8.8", "53")
	if err != nil {
		log.Println("Something went wrong during the lookup...")
		return
	}
	respBuffer := utils.PacketToBuf(*respPacket)
	socket.WriteTo(utils.From_512buffer(respBuffer.Buf), addr)
}
