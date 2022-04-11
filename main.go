package main

// TODO: Write REcord Config for IPv6 for AAAA query type

import (
	"github.com/orangeseeds/DNSserver/server"
	// "reflect"
)

func main() {
	udpServer := new(server.UPDServer)
	udpServer.SetHandler(server.HandleConnection)
	udpServer.Serve("1053")
}

// func format_packet(packet core.DnsPacket) {
// 	json_str, _ := utils.PrettyStruct(packet.Header)
// 	fmt.Println("header", json_str)
// 	for _, question := range packet.Questions {
// 		json_str, _ := utils.PrettyStruct(question)
// 		fmt.Println("question", json_str)
// 	}
// 	for _, answer := range packet.Answers {
// 		json_str, _ := utils.PrettyStruct(answer)
// 		fmt.Println("answer", json_str)
// 	}
// }
