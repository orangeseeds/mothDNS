package core_test

import (
	"os"
	"testing"

	"github.com/orangeseeds/DNSserver/core"
	"github.com/orangeseeds/DNSserver/utils"
)

func TestPacketConstr(t *testing.T) {

	f, _ := os.Open("response.txt")
	defer f.Close()

	buffer := core.NewBuffer()
	bef := make([]byte, 512)
	_, _ = f.Read(bef)

	buffer.Buf = bef

	d := new(core.Packet)
	e := *new(core.Record)
	var err error

	if err = d.Header.Read(&buffer); err != nil {
		t.Fatalf("error")
	}

	for i := 0; i < int(d.Header.Questions); i++ {

		question := core.NewQuestion("", core.QT_UNKNOWN)
		if err = question.Read(&buffer); err != nil {
			t.Fatalf("error")
		}

		d.Questions = append(d.Questions, question)
	}
	for i := 0; i < int(d.Header.Answers); i++ {

		if e, err = core.ReadRecord(&buffer); err != nil {
			t.Fatalf("error")
		}
		d.Answers = append(d.Answers, e)
	}

	t.Log(utils.PrettyStruct(d))
}
