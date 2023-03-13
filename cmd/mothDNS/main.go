package main

// TODO: Write REcord Config for IPv6 for AAAA query type

import (
	"github.com/orangeseeds/mothDNS/pkg/server"
)

func main() {
	udpServer := new(server.UPDServer)
	udpServer.SetHandler(server.HandleConnection)
	udpServer.Serve("1053")
	// apple
}
