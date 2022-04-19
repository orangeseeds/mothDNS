package server

import (
	"fmt"
	"log"
	"net"

	"github.com/orangeseeds/DNSserver/core"
)

// A helper function for setting some constraints on questions, doesn't really do much right now.
func CheckQuestions(questions []core.DnsQuestion) bool {
	for _, question := range questions {
		if question.Name == "127.0.0.1" {
			log.Printf("Question asking for %v, is not a valid question.", question.Name)
			return false
		}
	}
	return true
}

/*
   Sends a DNS request to a specified server

   @param name		-> domain name to lookup
   @param qType		-> DNS question type for the request
   @param serverType	-> server type of the target !Currently only "upd" works
   @param host		-> address of the server
   @param port		-> port number to send the packet to

   @return replyPacket	-> relpy from the server
*/
func LookUp(name string, qType core.QueryType, serverType string, host string, port string) (*core.Packet, error) {
	socket, err := net.Dial(serverType, host+":"+port)
	if err != nil {
		return nil, err
	}
	defer socket.Close()

	question := map[string]core.QueryType{
		name: qType,
	}

	packet := core.MakePacket(6666, true, question)
	buffer, _ := core.BuffFrmPacket(packet)

	_, err = socket.Write(buffer.Buf)
	if err != nil {
		return nil, fmt.Errorf("error while writing to %v, %v", host, err)
	}

	replyBuffer := make([]byte, 1000)
	_, err = socket.Read(replyBuffer)
	if err != nil {
		return nil, err
	}

	packetBuffer := core.NewBuffer()
	packetBuffer.Buf = replyBuffer
	replyPacket, err := core.PacketFrmBuff(&packetBuffer)
	if err != nil {
		return nil, err
	}

	return replyPacket, nil

}

/*
   Send recursive DNS requests to servers of different hierarchy.

   @param qname -> domain name to lookup
   @param qtype -> DNS question type for the request

   How does it works?
    - First we set a root server, in our case we chose a.root-server.net -> 198.41.0.4 and send out query to the server.
    - After querying the root server for our domain we get a list of TLDS(Top Level Domain Server) for our specific domain type .com/.net/.np/...
    - We then query the first TLDS that meets our requirement and get a list of name servers for the Domain
    - Any of the name servers should be able to resolve our Domain Name to IPadress
*/
func RecrLookUp(qname string, qtype core.QueryType) (*core.Packet, error) {

	ns := "198.41.0.4"

	for {
		log.Printf("looking up %v type:%v -> %v\n", qname, core.QtName(qtype), ns)
		response, _ := LookUp(qname, qtype, "udp", ns, "53")

		if len(response.Answers) != 0 && response.Header.Rescode == core.NOERROR {
			return response, nil
		}

		if response.Header.Rescode == core.NXDOMAIN {
			return response, nil
		}

		if newNs, err := response.GetResolvedNS(qname); err == nil {
			ns = newNs
			continue
		}

		var newName string
		if newNsName, err := response.GetUnresNS(qname); err != nil {
			return response, nil
		} else {
			newName = newNsName
		}

		recrResponse, _ := RecrLookUp(newName, core.QT_A)
		if newNs, err := recrResponse.GetRandomA(); err == nil {
			ns = string(newNs.(core.A).Addr.String())
		} else {
			return response, nil
		}
	}
}
