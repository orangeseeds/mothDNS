package core_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/orangeseeds/DNSserver/core"
	"github.com/orangeseeds/DNSserver/utils"
)

func TestPacket(t *testing.T) {

	f, _ := os.Open("response.txt")
	defer f.Close()

	buffer := core.NewBuffer()
	buff := make([]byte, 511, 512)
	_, _ = f.Read(buff)
	buffer.Buf = buff

	packet, _ := core.BufToPacket(buffer)

	fmt.Println(utils.PrettyStruct(packet))
}
