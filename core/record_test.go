package core_test

import (
	"os"
	"testing"

	"github.com/orangeseeds/DNSserver/core"
	"github.com/orangeseeds/DNSserver/utils"
)

func TestRecordBuilding(t *testing.T) {
	f, _ := os.Open("response_packet.txt")
	defer f.Close()

	buffer := core.NewBuffer()
	buff := make([]byte, 512)

	_, _ = f.Read(buff)
	buffer.Buf = buff

	t.Log(buffer)
	// t.Log(buffer)

	packet, _ := core.PacketFrmBuff(&buffer)
	ansBuf, _ := core.BuffFrmPacket(*packet)

	t.Log(ansBuf)
	// for _, r := range packet.Resources {
	// 	t.Logf("%T", r)
	// }
	t.Log(utils.PrettyStruct(packet))
	// respBuff, _ := core.BuffFrmPacket(*packet)
	// t.Log(respBuff.Buf[0:int(respBuff.Pos())])
}
