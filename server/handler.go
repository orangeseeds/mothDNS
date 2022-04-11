package server

import (
	"log"
	"net"

	"github.com/orangeseeds/DNSserver/core"
	"github.com/orangeseeds/DNSserver/utils"
)

func HandleConnection(socket net.PacketConn, addr net.Addr, buf []byte) {
	rcvBuffer := core.NewBuffer()
	rcvBuffer.Buf = utils.To_512buffer(buf)
	rcvPacket, err := utils.BufToPacket(rcvBuffer)
	if err != nil {
		log.Println("request buffer to packet conversion failed...")
		return
	}
	if !utils.CheckQuestions(rcvPacket.Questions) {
		socket.WriteTo(utils.From_512buffer(rcvBuffer.Buf), addr)
		return
	}
	nameQtypes := make(map[string]core.QueryType)
	for _, question := range rcvPacket.Questions {
		nameQtypes[question.Name] = question.Qtype
	}

	host := "8.8.8.8"
	respPacket, err := utils.Lookup(nameQtypes, rcvPacket.Header.Id, "udp", host, "53")
	if err != nil {
		log.Println("Something went wrong during the lookup...")
		return
	}
	log.Printf(`%v -> %v | successfully resolved!`, nameQtypes, addr)
	respBuffer := utils.PacketToBuf(*respPacket)
	socket.WriteTo(utils.From_512buffer(respBuffer.Buf), addr)
}
