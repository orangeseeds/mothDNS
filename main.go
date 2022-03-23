package main

import (
	"DNSserver/packages"
	"fmt"
	"os"
	// "reflect"
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

	qname := "google.com"
	qtype := header.qt_A
}

func main() {
	// fmt.Println("")
	// read_test()
	stub_resolver()
}
