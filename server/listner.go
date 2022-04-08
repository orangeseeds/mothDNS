package main

import (
	"fmt"
	"log"
	"net"
)

func serve(socket net.PacketConn, addr net.Addr, buf []byte) {

	message := "Heyy there! Client..."
	fmt.Println("Something just hit me!!", string(buf))

	socket.WriteTo([]byte(message), addr)
}

func main() {
	// listen to incoming udp packets
	socket, err := net.ListenPacket("udp", ":1053")
	if err != nil {
		log.Fatal(err)
	}
	defer socket.Close()

	for {
		buf := make([]byte, 512)
		_, addr, err := socket.ReadFrom(buf)
		if err != nil {
			continue
		}
		go serve(socket, addr, buf)
	}

}
