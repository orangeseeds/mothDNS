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

func main() {
	fmt.Println("")
	f, err := os.Open("response_packet.txt")
	defer f.Close()
	check(err)

	var buffer = header.NewBuffer()
	fmt.Println(buffer.Pos())
	buff := make([]byte, 512)
	var num, _ = f.Read(buff)

	// fmt.Println(buff)
	// fmt.Println(num)
	for i := range buff {
		buffer.Buf[i] = buff[i]
	}
	//
	// buffer.Buf[0] = 134
	// buffer.Buf[1] = 42
	//
	packet := header.NewPacket()
	var _, _ = packet.From_buffer(&buffer)
	// packet = *p_packet
	//
	// // fmt.Println(res)
	// // fmt.Println("assss")
	// // fmt.Println(num)
	fmt.Println(num)
}
