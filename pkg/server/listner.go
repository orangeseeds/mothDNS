package server

import (
	"log"
	"net"
)

type RqstHandler func(net.PacketConn, net.Addr, []byte, string)

type UPDServer struct {
	RootServer string
	Handler    RqstHandler
}

func (s *UPDServer) Serve(port string) *UPDServer {

	socket, err := net.ListenPacket("udp", port)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Listening on port " + port + "...")
	defer socket.Close()

	for {
		buf := make([]byte, 512)
		_, addr, err := socket.ReadFrom(buf)
		if err != nil {
			log.Println(err)
			continue
		}
		go s.Handler(socket, addr, buf, s.RootServer)
	}

}

func (s *UPDServer) SetHandler(f RqstHandler) {
	s.Handler = f
}
