package core

import (
	"github.com/orangeseeds/mothDNS/pkg/bpb"
)

type DNSQuestion struct {
	Name  string
	Qtype QueryType
}

func NewQuestion(name string, qtype QueryType) DNSQuestion {
	q := DNSQuestion{
		Name:  name,
		Qtype: qtype,
	}
	return q
}

func (q *DNSQuestion) Read(buffer *bpb.BytePacketBuffer) error {

	var err error
	buffer.ReadQName(&q.Name)

	var result uint16
	if result, err = buffer.ReadTwoBytes(); err != nil {
		return err
	}
	r_buffer := result

	q.Qtype = QueryType.From_num(0, r_buffer)

	if result, err = buffer.ReadTwoBytes(); err != nil {
		return err
	}
	var _ = result

	return nil
}

// ---------------------------------- For Writing ---------------------------------------------------

func (d *DNSQuestion) Write(buffer *bpb.BytePacketBuffer) error {
	if err := buffer.WriteQName(&d.Name); err != nil {
		return err
	}

	typenum := uint16(d.Qtype)
	if err := buffer.WriteTwoBytes(typenum); err != nil {
		return err
	}
	if err := buffer.WriteTwoBytes(1); err != nil {
		return err
	}

	return nil
}
