package server

import (
	"fmt"
	"log"
	"net"
)

type RqstHandler func(net.PacketConn, net.Addr, []byte)

type UPDServer struct {
	Handler RqstHandler
}

func (s *UPDServer) Serve(port string) *UPDServer {

	socket, err := net.ListenPacket("udp", ":"+port)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Listening on port:" + port + "...")
	defer socket.Close()

	for {
		buf := make([]byte, 512)
		_, addr, err := socket.ReadFrom(buf)
		if err != nil {
			log.Println(err)
			continue
		}
		go s.Handler(socket, addr, buf)
	}

}

func (s *UPDServer) SetHandler(f RqstHandler) {
	s.Handler = f
}

func serve(socket net.PacketConn, addr net.Addr, buf []byte) {
	message := "Heyy there! Client..."
	fmt.Println("Something just hit me!!", string(buf))

	socket.WriteTo([]byte(message), addr)
}

// func main() {
// 	server := new(UPDServer)
// 	server.SetHandler(serve)
// 	server.Serve("1053")
// }
