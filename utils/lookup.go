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

func LookUp(name string, qType core.QueryType, serverType string, host string, port string) (*core.DnsPacket, error) {
	socket, err := net.Dial(serverType, host+":"+port)
	if err != nil {
		return nil, err
	}
	defer socket.Close()

	question := map[string]core.QueryType{
		name: qType,
	}

	fmt.Println(question)

	packet := ConstrPacket(6666, true, question)
	buffer := PacketToBuf(packet)

	fmt.Println(packet.Questions)
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

func RecrLookUp(qname string, qtype core.QueryType) (*core.DnsPacket, error) {

	ns := "193.0.14.129"

	for {
		fmt.Printf("attempting lookup for %v type:%v -> %v\n", qname, qtype, ns)

		response, err := LookUp(qname, qtype, "udp", ns, "53")
		fmt.Printf("Status -> %v", response.Header.Rescode)
		if err != nil {
			return nil, err
		}

		// val, _ := json.MarshalIndent(response, "", "    ")
		// fmt.Println("sallaaa", string(val))

		if len(response.Answers) != 0 && response.Header.Rescode == core.NOERROR {
			return response, nil
		}

		if response.Header.Rescode == core.NXDOMAIN {
			return response, nil
		}

		// val, _ := json.MarshalIndent(response, "", "    ")
		// fmt.Println(string(val))
		if newNs, err := response.GetResolvedNS(qname); err == nil {
			ns = newNs
			fmt.Println("\n new:" + newNs)
			continue
		} else {
			return response, nil
		}

		// return response, nil
	}
}
