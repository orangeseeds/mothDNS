package core

// Some utility function to reuse in the code
func MakePacket(id uint16, isRec bool, nameQtypes map[string]QueryType) Packet {
	packet := Packet{}
	packet.Header.Id = id
	packet.Header.Recursion_desired = isRec

	for name, qtype := range nameQtypes {
		packet.Questions = append(packet.Questions, NewQuestion(name, qtype))
	}
	packet.Header.Questions = uint16(len(nameQtypes))
	return packet
}
