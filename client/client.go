package main

import (
	// "bufio"

	"encoding/json"
	"fmt"
	"net"

	"github.com/orangeseeds/mothDNS/pkg/bpb"
	"github.com/orangeseeds/mothDNS/pkg/core"
)

// func main() {
// 	p := make([]byte, 512)
// 	conn, err := net.Dial("udp", "127.0.0.1:1053")
// 	if err != nil {
// 		fmt.Printf("Some error %v", err)
// 		return
// 	}
// 	defer conn.Close()

// 	_, err = conn.Write([]byte("Hello UDP server!"))
// 	_, err = conn.Read(p)
// 	if err == nil {
// 		fmt.Println(string(p))
// 	} else {
// 		fmt.Printf("Some error %v\n", err)
// 	}
// }

const (
	SERVER_HOST = "192.12.94.30"
	SERVER_PORT = "53"
	SERVER_TYPE = "udp"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	name := "yahoo.com"
	packet, _ := lookup(name, core.QT_A)

	val, _ := json.MarshalIndent(packet, "", "    ")
	fmt.Println(string(val))
}

func lookup(qname string, qtype core.QueryType) (core.DNSPacket, error) {
	socket, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	check(err)

	defer socket.Close()

	packet := core.NewPacket()

	packet.Header.Id = 6666
	packet.Header.Questions = 1
	packet.Header.RecursionDesired = true
	packet.Questions = append(packet.Questions, core.NewQuestion(string(qname), qtype))

	req_buffer := bpb.New()
	packet.Write(&req_buffer)

	buff := make([]byte, 250)
	for i := range buff {
		buff[i] = req_buffer.Buf[i]
	}

	_, err = socket.Write(buff)

	check(err)

	reply := make([]byte, 512)

	_, _ = socket.Read(reply)
	r_buffer := bpb.New()

	r_buffer.Buf = reply
	// fmt.Println(reply_buff)

	res_packet := core.NewPacket()
	res_packet.From_buffer(&r_buffer)
	return res_packet, nil
}
