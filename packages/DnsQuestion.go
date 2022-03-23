package header

import (
// "errors"
// "fmt"
// "strings"
)

type DnsQuestion struct {
	name  string
	qtype QueryType
}

func NewQuestion(name string, qtype QueryType) DnsQuestion {
	q := DnsQuestion{
		name:  name,
		qtype: qtype,
	}
	return q
}

func (q *DnsQuestion) Read(buffer *BytePacketBuffer) error {

	var err error
	buffer.Read_qname(&q.name)

	var result uint16
	if result, err = buffer.Read_u16(); err != nil {
		return err
	}
	r_buffer := result

	q.qtype = QueryType.From_num(0, r_buffer)

	if result, err = buffer.Read_u16(); err != nil {
		return err
	}
	var _ = result

	return nil
}

// ################################## For Writing ###################################################

func (d *DnsQuestion) Write(buffer *BytePacketBuffer) error {
	if err := buffer.Write_qname(&d.name); err != nil {
		return err
	}

	typenum := uint16(d.qtype)
	if err := buffer.Write_u16(typenum); err != nil {
		return err
	}
	if err := buffer.Write_u16(1); err != nil {
		return err
	}

	return nil
}
