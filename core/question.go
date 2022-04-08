package core

import (
// "errors"
// "fmt"
// "strings"
)

type DnsQuestion struct {
	Name  string
	Qtype QueryType
}

func NewQuestion(name string, qtype QueryType) DnsQuestion {
	q := DnsQuestion{
		Name:  name,
		Qtype: qtype,
	}
	return q
}

func (q *DnsQuestion) Read(buffer *BytePacketBuffer) error {

	var err error
	buffer.Read_qname(&q.Name)

	var result uint16
	if result, err = buffer.Read_u16(); err != nil {
		return err
	}
	r_buffer := result

	q.Qtype = QueryType.From_num(0, r_buffer)

	if result, err = buffer.Read_u16(); err != nil {
		return err
	}
	var _ = result

	return nil
}

// ################################## For Writing ###################################################

func (d *DnsQuestion) Write(buffer *BytePacketBuffer) error {
	if err := buffer.Write_qname(&d.Name); err != nil {
		return err
	}

	typenum := uint16(d.Qtype)
	if err := buffer.Write_u16(typenum); err != nil {
		return err
	}
	if err := buffer.Write_u16(1); err != nil {
		return err
	}

	return nil
}
