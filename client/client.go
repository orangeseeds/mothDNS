package main

import (
	// "bufio"
	"fmt"
	"net"
)

func main() {
	p := make([]byte, 512)
	conn, err := net.Dial("udp", "127.0.0.1:1053")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte("Hello UDP server!"))
	_, err = conn.Read(p)
	if err == nil {
		fmt.Println(string(p))
	} else {
		fmt.Printf("Some error %v\n", err)
	}
}
