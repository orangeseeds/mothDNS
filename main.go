package main

// TODO: Write REcord Config for IPv6 for AAAA query type

import (
	"DNSserver/packages"
	"fmt"
	"net"
	"os"
	// "reflect"
)

const (
	SERVER_HOST = "8.8.8.8"
	SERVER_PORT = "53"
	SERVER_TYPE = "udp"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func read_test() {
	f, err := os.Open("response_packet.txt")
	defer f.Close()
	check(err)

	var buffer = header.NewBuffer()
	buff := make([]byte, 512)
	var _, _ = f.Read(buff)

	for i := range buff {
		buffer.Buf[i] = buff[i]
	}
	packet := header.NewPacket()
	var _, _ = packet.From_buffer(&buffer)
	fmt.Printf("%T: %+v\n", packet.Header, packet.Header)
	fmt.Printf("%T: %+v\n", packet.Questions, packet.Questions)
	fmt.Printf("%T: %+v\n", packet.Answers, packet.Answers)
}

func stub_resolver() {

	qname := "yahoo.com"
	qtype := header.QT_MX

	socket, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	check(err)

	defer socket.Close()

	packet := header.NewPacket()

	packet.Header.Id = 6666
	packet.Header.Questions = 1
	packet.Header.Recursion_desired = true
	packet.Questions = append(packet.Questions, header.NewQuestion(string(qname), qtype))

	req_buffer := header.NewBuffer()
	packet.Write(&req_buffer)

	buff := make([]byte, 512)
	for i := range buff {
		buff[i] = req_buffer.Buf[i]
	}

	_, err = socket.Write(buff)

	check(err)

	reply := make([]byte, 512)

	_, err = socket.Read(reply)

	r_buffer := header.NewBuffer()
	reply_buff := [512]byte{0}
	for i := range buff {
		reply_buff[i] = reply[i]
	}

	r_buffer.Buf = reply_buff

	res_packet := header.NewPacket()
	res_packet.From_buffer(&r_buffer)
	fmt.Printf("%T: %+v\n", res_packet.Header, res_packet.Header)
	fmt.Printf("%T: %+v\n", res_packet.Questions, res_packet.Questions)
	fmt.Printf("%T: %+v\n", res_packet.Answers, res_packet.Answers)
	// fmt.Printf("%T: %+v\n", res_packet.Answers, res_packet.Answers)

}

func net_test() {
	if len(os.Args) != 2 {

		fmt.Println("Usage: echo_client message")
		os.Exit(1)
	}

	msg := os.Args[1]

	con, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)

	check(err)

	defer con.Close()

	_, err = con.Write([]byte(msg))

	check(err)

	reply := make([]byte, 1024)

	_, err = con.Read(reply)

	check(err)

	fmt.Println("reply:", string(reply))
}

func lookup(qname string, qtype header.QueryType) (header.DnsPacket, error) {
	socket, err := net.Dial(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	check(err)

	defer socket.Close()

	packet := header.NewPacket()

	packet.Header.Id = 6666
	packet.Header.Questions = 1
	packet.Header.Recursion_desired = true
	packet.Questions = append(packet.Questions, header.NewQuestion(string(qname), qtype))

	req_buffer := header.NewBuffer()
	packet.Write(&req_buffer)

	buff := make([]byte, 512)
	for i := range buff {
		buff[i] = req_buffer.Buf[i]
	}

	_, err = socket.Write(buff)

	check(err)

	reply := make([]byte, 512)

	_, err = socket.Read(reply)

	r_buffer := header.NewBuffer()
	reply_buff := [512]byte{0}
	for i := range buff {
		reply_buff[i] = reply[i]
	}

	r_buffer.Buf = reply_buff

	res_packet := header.NewPacket()
	res_packet.From_buffer(&r_buffer)
	return res_packet, nil
}

func main() {
	// fmt.Println("")
	// read_test()
	name := "google.com"
	packet, _ := lookup(name, header.QT_A)

	fmt.Printf("%T: %+v\n", packet.Header, packet.Header)
	fmt.Printf("%T: %+v\n", packet.Questions, packet.Questions)
	fmt.Printf("%T: %+v\n", packet.Answers, packet.Answers)
	// net_test()
}

// server, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
// if err != nil {
// 	fmt.Println("Error listening:", err.Error())
// 	os.Exit(1)
// }
//
// defer server.Close()
// fmt.Println("Listening on " + SERVER_HOST + ":" + SERVER_PORT)
// fmt.Println("Waiting for client...")
