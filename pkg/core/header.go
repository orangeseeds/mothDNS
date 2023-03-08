package core

import (
	"github.com/orangeseeds/mothDNS/pkg/bpb"
)

type dnsHeader struct {
	Id                   uint16     `json:"id"`                    // 16 bits
	Response             bool       `json:"response"`              // 1 bit
	Opcode               byte       `json:"opcode"`                // 4 bits
	AuthoritativeAnswer  bool       `json:"authoritative_answer"`  // 1 bit
	TruncatedMessage     bool       `json:"truncated_message"`     // 1 bit
	RecursionDesired     bool       `json:"recursion_desired"`     // 1 bit
	RecursionAvailable   bool       `json:"recursion_available"`   // 1 bit
	Z                    bool       `json:"z"`                     // 1 bit
	Rescode              ResultCode `json:"rescode"`               // 4 bits
	CheckingDisabled     bool       `json:"checking_disabled"`     // 1 bit
	Questions            uint16     `json:"questions"`             // 16 bits
	Answers              uint16     `json:"answers"`               // 16 bits
	AuthoritativeEntries uint16     `json:"authoritative_entries"` // 16 bits
	AuthedData           bool       `json:"authed_data"`           // 1 bit
	ResourceEntries      uint16     `json:"resource_entries"`      // 16 bits
}

func NewHeader() *dnsHeader {
	d := dnsHeader{
		Rescode: NOERROR,
	}
	return &d
}

func (d *dnsHeader) read(buffer *bpb.BytePacketBuffer) error {

	return nil
}

func (d *dnsHeader) Read(buffer *bpb.BytePacketBuffer) error {
	var (
		err      error
		twoBytes uint16
	)

	// reading ID from first 16 bits of buffer
	if twoBytes, err = buffer.ReadTwoBytes(); err != nil {
		return err
	}
	d.Id = twoBytes

	{
		// reading next 16 bits from QR to RCODE
		if twoBytes, err = buffer.ReadTwoBytes(); err != nil {
			return err
		}
		flags := twoBytes

		a := uint8((flags >> 8))
		b := uint8((flags & 0xFF))

		d.RecursionDesired = (a & (1 << 0)) > 0
		d.TruncatedMessage = (a & (1 << 1)) > 0
		d.AuthoritativeAnswer = (a & (1 << 2)) > 0
		d.Opcode = (a >> 3) & 0x0F
		d.Response = (a & (1 << 7)) > 0

		d.Rescode = ResultCode.From_num(0, b&0x0F)
		d.CheckingDisabled = (b & (1 << 4)) > 0
		d.AuthedData = (b & (1 << 5)) > 0
		d.Z = (b & (1 << 6)) > 0
		d.RecursionAvailable = (b & (1 << 7)) > 0
	}

	{
		// reading question count from next 16 bits
		if twoBytes, err = buffer.ReadTwoBytes(); err != nil {
			return err
		}
		d.Questions = twoBytes
	}

	{
		// reading answer count from next 16 bits
		if twoBytes, err = buffer.ReadTwoBytes(); err != nil {
			return err
		}
		d.Answers = twoBytes
	}

	{
		// reading Authority count from next 16 bits
		if twoBytes, err = buffer.ReadTwoBytes(); err != nil {
			return err
		}
		d.AuthoritativeEntries = twoBytes
	}

	{
		// reading resource entries from next 16 bits
		if twoBytes, err = buffer.ReadTwoBytes(); err != nil {
			return err
		}
		d.ResourceEntries = twoBytes
	}

	// fmt.Printf("%+v\n", d)
	return nil
}

// ---------------------------------- For Writing ---------------------------------------------------

func To_uint8(v bool) uint8 {
	num_u8 := uint8(0)
	if v {
		num_u8 = 1
	}

	return num_u8
}

func (d *dnsHeader) Write(buffer *bpb.BytePacketBuffer) error {

	if err := buffer.WriteTwoBytes(d.Id); err != nil {
		return err
	}

	v1 := To_uint8(d.RecursionDesired)
	v2 := (To_uint8(d.TruncatedMessage) << 1)
	v3 := (To_uint8(d.AuthoritativeAnswer) << 2)
	v4 := (d.Opcode << 3)
	v5 := uint8(To_uint8(d.Response) << 7)

	if err := buffer.WriteOneByte(v1 | v2 | v3 | v4 | v5); err != nil {
		return err
	}

	v1 = uint8(d.Rescode)
	v2 = To_uint8(d.CheckingDisabled) << 4
	v3 = To_uint8(d.AuthedData) << 5
	v4 = To_uint8(d.Z) << 6
	v5 = To_uint8(d.RecursionAvailable) << 7

	if err := buffer.WriteOneByte(v1 | v2 | v3 | v4 | v5); err != nil {
		return err
	}

	if err := buffer.WriteTwoBytes(d.Questions); err != nil {
		return err
	}
	if err := buffer.WriteTwoBytes(d.Answers); err != nil {
		return err
	}
	if err := buffer.WriteTwoBytes(d.AuthoritativeEntries); err != nil {
		return err
	}
	if err := buffer.WriteTwoBytes(d.ResourceEntries); err != nil {
		return err
	}

	return nil
}
