package core

import ()

type DnsHeader struct {
	Id                    uint16     `json:"id"`                    // 16 bits
	Recursion_desired     bool       `json:"recursion_desired"`     // 1 bit
	Truncated_message     bool       `json:"truncated_message"`     // 1 bit
	Authoritative_answer  bool       `json:"authoritative_answer"`  // 1 bit
	Opcode                byte       `json:"opcode"`                // 4 bits
	Response              bool       `json:"response"`              // 1 bit
	Rescode               ResultCode `json:"rescode"`               // 4 bits
	Checking_disabled     bool       `json:"checking_disabled"`     // 1 bit
	Authed_data           bool       `json:"authed_data"`           // 1 bit
	Z                     bool       `json:"z"`                     // 1 bit
	Recursion_available   bool       `json:"recursion_available"`   // 1 bit
	Questions             uint16     `json:"questions"`             // 16 bits
	Answers               uint16     `json:"answers"`               // 16 bits
	Authoritative_entries uint16     `json:"authoritative_entries"` // 16 bits
	Resource_entries      uint16     `json:"resource_entries"`      // 16 bits
}

func NewHeader() *DnsHeader {
	var rCode ResultCode
	rCode = NOERROR

	d := DnsHeader{
		Id:                    0,
		Recursion_desired:     false,
		Truncated_message:     false,
		Authoritative_answer:  false,
		Opcode:                0,
		Response:              false,
		Rescode:               rCode,
		Checking_disabled:     false,
		Authed_data:           false,
		Z:                     false,
		Recursion_available:   false,
		Questions:             0,
		Answers:               0,
		Authoritative_entries: 0,
		Resource_entries:      0,
	}
	return &d
}

func (d *DnsHeader) Read(buffer *BytePacketBuffer) error {
	var err error

	var byte2 uint16
	if byte2, err = buffer.Read_u16(); err != nil {
		return err
	}
	d.Id = byte2

	if byte2, err = buffer.Read_u16(); err != nil {
		return err
	}
	flags := byte2

	a := uint8((flags >> 8))
	b := uint8((flags & 0xFF))

	d.Recursion_desired = (a & (1 << 0)) > 0
	d.Truncated_message = (a & (1 << 1)) > 0
	d.Authoritative_answer = (a & (1 << 2)) > 0
	d.Opcode = (a >> 3) & 0x0F
	d.Response = (a & (1 << 7)) > 0

	d.Rescode = ResultCode.From_num(0, b&0x0F)
	d.Checking_disabled = (b & (1 << 4)) > 0
	d.Authed_data = (b & (1 << 5)) > 0
	d.Z = (b & (1 << 6)) > 0
	d.Recursion_available = (b & (1 << 7)) > 0

	if byte2, err = buffer.Read_u16(); err != nil {
		return err
	}
	d.Questions = byte2

	if byte2, err = buffer.Read_u16(); err != nil {
		return err
	}
	d.Answers = byte2

	if byte2, err = buffer.Read_u16(); err != nil {
		return err
	}
	d.Authoritative_entries = byte2

	if byte2, err = buffer.Read_u16(); err != nil {
		return err
	}
	d.Resource_entries = byte2

	return nil
}

// ################################## For Writing ###################################################

func To_uint8(v bool) uint8 {
	num_u8 := uint8(0)
	if v {
		num_u8 = 1
	}

	return num_u8
}

func (d *DnsHeader) Write(buffer *BytePacketBuffer) error {

	if err := buffer.Write_u16(d.Id); err != nil {
		return err
	}

	v1 := To_uint8(d.Recursion_desired)
	v2 := (To_uint8(d.Truncated_message) << 1)
	v3 := (To_uint8(d.Authoritative_answer) << 2)
	v4 := (d.Opcode << 3)
	v5 := uint8(To_uint8(d.Response) << 7)

	if err := buffer.Write_u8(v1 | v2 | v3 | v4 | v5); err != nil {
		return err
	}

	v1 = uint8(d.Rescode)
	v2 = To_uint8(d.Checking_disabled) << 4
	v3 = To_uint8(d.Authed_data) << 5
	v4 = To_uint8(d.Z) << 6
	v5 = To_uint8(d.Recursion_available) << 7

	if err := buffer.Write_u8(v1 | v2 | v3 | v4 | v5); err != nil {
		return err
	}

	if err := buffer.Write_u16(d.Questions); err != nil {
		return err
	}
	if err := buffer.Write_u16(d.Answers); err != nil {
		return err
	}
	if err := buffer.Write_u16(d.Authoritative_entries); err != nil {
		return err
	}
	if err := buffer.Write_u16(d.Resource_entries); err != nil {
		return err
	}

	return nil
}
