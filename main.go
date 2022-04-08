package main

// TODO: Write REcord Config for IPv6 for AAAA query type

import (
	"DNSserver/core"
	"DNSserver/utils"
	"fmt"
	"log"
	"net"
	// "reflect"
)

const (
	SERVER_HOST = "8.8.8.8"
	SERVER_PORT = "53"
	SERVER_TYPE = "udp"
)

func main() {

	// name := "google.com"
	// packet, _ := lookup(name, core.QT_A)
	// format_packet(packet)

	socket, err := net.ListenPacket("udp", "127.0.0.1:2053")
	if err != nil {
		log.Fatal(err)
	}
	// defer socket.Close()

	for {
		// fmt.Println("bruvv")
		buf := make([]byte, 512)
		_, addr, err := socket.ReadFrom(buf)
		if err != nil {
			continue
		}
		go handleConnection(socket, addr, buf)
	}
}

func handleConnection(socket net.PacketConn, addr net.Addr, buf []byte) {

	r_buffer := core.NewBuffer()
	r_buffer.Buf = utils.To_512buffer(buf)

	req_packet := core.NewPacket()
	req_packet.From_buffer(&r_buffer)

	// fmt.Println(req_packet)

	// format_packet(req_packet)
	name := req_packet.Questions[0].Name
	id := req_packet.Header.Id

	if name == "127.0.0.1" {
		socket.WriteTo(utils.From_512buffer(r_buffer.Buf), addr)
		return
	}
	// fmt.Println(name)

	// request_packet := utils.From_512buffer(r_buffer.Buf)
	// fmt.Println("Something just hit me!!", string(buf))

	// name := "google.com"
	packet, _ := lookup(name, core.QT_A, id)
	format_packet(packet)
	res_buffer := core.NewBuffer()
	packet.Write(&res_buffer)

	socket.WriteTo(utils.From_512buffer(res_buffer.Buf), addr)
}

func lookup(qname string, qtype core.QueryType, id uint16) (core.DnsPacket, error) {
	socket, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		panic(err)
	}

	defer socket.Close()

	packet := core.NewPacket()
	packet.Header.Id = id
	packet.Header.Questions = 1
	packet.Header.Recursion_desired = true
	packet.Questions = append(packet.Questions, core.NewQuestion(string(qname), qtype))

	req_buffer := core.NewBuffer()
	packet.Write(&req_buffer)

	buff := utils.From_512buffer(req_buffer.Buf)

	_, err = socket.Write(buff)
	check(err)

	reply := make([]byte, 512)
	_, err = socket.Read(reply)
	r_buffer := core.NewBuffer()
	r_buffer.Buf = utils.To_512buffer(reply)

	res_packet := core.NewPacket()
	res_packet.From_buffer(&r_buffer)
	return res_packet, nil
}

func format_packet(packet core.DnsPacket) {
	json_str, _ := utils.PrettyStruct(packet.Header)
	fmt.Println("header", json_str)
	for _, question := range packet.Questions {
		json_str, _ := utils.PrettyStruct(question)
		fmt.Println("question", json_str)
	}
	for _, answer := range packet.Answers {
		json_str, _ := utils.PrettyStruct(answer)
		fmt.Println("answer", json_str)
	}
}
