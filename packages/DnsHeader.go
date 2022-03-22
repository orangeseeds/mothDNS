package header

import (
	"errors"
	"fmt"
	// "strings"
)

var EOBError = errors.New("Buffer greater than or equals 512")

type DnsHeader struct {
	id                    uint16     // 16 bits
	recursion_desired     bool       // 1 bit
	truncated_message     bool       // 1 bit
	authoritative_answer  bool       // 1 bit
	opcode                byte       // 4 bits
	response              bool       // 1 bit
	rescode               ResultCode // 4 bits
	checking_disabled     bool       // 1 bit
	authed_data           bool       // 1 bit
	z                     bool       // 1 bit
	recursion_available   bool       // 1 bit
	questions             uint16     // 16 bits
	answers               uint16     // 16 bits
	authoritative_entries uint16     // 16 bits
	resource_entries      uint16     // 16 bits
}

func NewHeader() DnsHeader {
	fmt.Println("")
	var rCode ResultCode
	rCode = NOERROR

	d := DnsHeader{
		id:                    0,
		recursion_desired:     false,
		truncated_message:     false,
		authoritative_answer:  false,
		opcode:                0,
		response:              false,
		rescode:               rCode,
		checking_disabled:     false,
		authed_data:           false,
		z:                     false,
		recursion_available:   false,
		questions:             0,
		answers:               0,
		authoritative_entries: 0,
		resource_entries:      0,
	}

	return d
}

func (d *DnsHeader) Read(buffer *BytePacketBuffer) error {
	// fmt.Println(buffer)
	var err error

	var byte2 *uint16
	if byte2, err = buffer.Read_u16(); err != nil {
		return err
	}
	d.id = *byte2

	if byte2, err = buffer.Read_u16(); err != nil {
		return err
	}
	flags := *byte2

	a := uint8((flags >> 8))
	b := uint8((flags & 0xFF))

	d.recursion_desired = (a & (1 << 0)) > 0
	d.truncated_message = (a & (1 << 1)) > 0
	d.authoritative_answer = (a & (1 << 2)) > 0
	d.opcode = (a >> 3) & 0x0F
	d.response = (a & (1 << 7)) > 0

	d.rescode = ResultCode.From_num(0, b&0x0F)
	d.checking_disabled = (b & (1 << 4)) > 0
	d.authed_data = (b & (1 << 5)) > 0
	d.z = (b & (1 << 6)) > 0
	d.recursion_available = (b & (1 << 7)) > 0

	if byte2, err = buffer.Read_u16(); err != nil {
		return err
	}
	d.questions = *byte2

	if byte2, err = buffer.Read_u16(); err != nil {
		return err
	}
	d.answers = *byte2

	if byte2, err = buffer.Read_u16(); err != nil {
		return err
	}
	d.authoritative_entries = *byte2

	if byte2, err = buffer.Read_u16(); err != nil {
		return err
	}
	d.resource_entries = *byte2

	return nil
}
