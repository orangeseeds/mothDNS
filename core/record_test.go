package core_test

import (
	"testing"

	"github.com/orangeseeds/DNSserver/core"
	"github.com/orangeseeds/DNSserver/utils"
)

func TestUnkRecord(t *testing.T) {
	recordU := core.UNKNOWN{
		Domain:   "google.com",
		Qtype:    32,
		Data_len: 600,
		Ttl:      800,
	}

	recordA := core.A{
		Domain: "yahoo.com",
		Ttl:    900,
	}

	records := []core.Record{}
	records = append(records, recordU)
	records = append(records, recordA)

	for _, r := range records {
		val, _ := utils.PrettyStruct(r)
		switch r.(type) {
		case core.UNKNOWN:
			t.Logf("UNKNOWN %v", val)

		case core.A:
			t.Logf("A %v", val)
		}
	}
}
