package server

import (
	"log"
	"net"

	"github.com/orangeseeds/DNSserver/core"
)

/*
   Handles the DNS question sent to the server.

   @param socket   -> generic packet oriented generic connection.
   @param addr	    -> address of the request sender
   @param buf	    -> byte buffer received from the request

*/
func HandleConnection(socket net.PacketConn, addr net.Addr, buf []byte) {
	rcvBuffer := core.NewBuffer()
	rcvBuffer.Buf = buf
	rcvPacket, err := core.BufToPacket(rcvBuffer)
	if err != nil {
		log.Println("request buffer to packet conversion failed...")
		return
	}
	if !CheckQuestions(rcvPacket.Questions) {
		socket.WriteTo(rcvBuffer.Buf, addr)
		return
	}

	replyPacket := core.NewPacket()
	replyPacket.Header.Id = rcvPacket.Header.Id
	replyPacket.Header.Recursion_desired = true
	replyPacket.Header.Recursion_available = true
	replyPacket.Header.Response = true

	if len(rcvPacket.Questions) > 0 {
		for _, q := range rcvPacket.Questions {

			respPacket, err := RecrLookUp(q.Name, q.Qtype)
			if err != nil {
				log.Println("Something went wrong during the lookup...")
				replyPacket.Header.Rescode = core.SERVFAIL
			} else {
				log.Printf(`%v -> %v type: %v | resolved!`, addr, q.Name, core.QtName(q.Qtype))

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
	replyBuffer := core.PacketToBuf(replyPacket)
	socket.WriteTo(replyBuffer.Buf, addr)

}
