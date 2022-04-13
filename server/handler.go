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

	// host := "8.8.8.8"

	replyPacket := core.NewPacket()
	replyPacket.Header.Id = rcvPacket.Header.Id
	replyPacket.Header.Recursion_desired = true
	replyPacket.Header.Recursion_available = true
	replyPacket.Header.Response = true

	if len(rcvPacket.Questions) > 0 {
		for _, q := range rcvPacket.Questions {

			respPacket, err := utils.RecrLookUp(q.Name, q.Qtype)
			if err != nil {
				log.Println("Something went wrong during the lookup...")
				replyPacket.Header.Rescode = core.SERVFAIL
			} else {
				log.Printf(`%v -> %v type: %v | successfully resolved!`, addr, q.Name, q.Qtype)

				replyPacket.Questions = append(replyPacket.Questions, q)
				replyPacket.Header.Rescode = respPacket.Header.Rescode

				replyPacket.Answers = respPacket.Answers
				replyPacket.Authorities = respPacket.Authorities

				replyPacket.Resources = respPacket.Resources

			}

		}
	} else {
		replyPacket.Header.Rescode = core.FORMERR
	}
	replyBuffer := utils.PacketToBuf(replyPacket)
	socket.WriteTo(utils.From_512buffer(replyBuffer.Buf), addr)

}
