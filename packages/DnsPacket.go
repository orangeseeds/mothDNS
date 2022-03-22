package header

import (
	"fmt"
)

type DnsPacket struct {
	Header      DnsHeader
	questions   []DnsQuestion
	answers     []DnsRecord
	authorities []DnsRecord
	resources   []DnsRecord
}

func NewPacket() DnsPacket {
	d := DnsPacket{
		Header:      *new(DnsHeader),
		questions:   []DnsQuestion{},
		answers:     []DnsRecord{},
		authorities: []DnsRecord{},
		resources:   []DnsRecord{},
	}

	return d
}

func (d *DnsPacket) From_buffer(buffer *BytePacketBuffer) (*DnsPacket, error) {
	// fmt.Println(buffer)
	var err error

	if err = d.Header.Read(buffer); err != nil {
		return nil, err
	}
	fmt.Println(d.Header)

	fmt.Println("DNS Questions")
	for i := 0; i < int(d.Header.questions); i++ {

		question := NewQuestion("", qt_UNKNOWN)
		if err = question.Read(buffer); err != nil {
			return nil, err
		}

		d.questions = append(d.questions, question)
		fmt.Println(buffer.pos)
	}
	fmt.Println(d.questions)
	//
	for i := 0; i < int(d.Header.answers); i++ {

		var p_rec *DnsRecord
		if p_rec, err = eDnsRecord.Read(buffer); err != nil {
			return nil, err
		}
		d.answers = append(d.answers, *p_rec)
		fmt.Println(buffer.pos)
	}
	fmt.Println(d.answers)
	//
	// for i := 0; i < int(d.Header.authoritative_entries); i++ {
	//
	// 	var p_rec *DnsRecord
	// 	if p_rec, err = eDnsRecord.Read(*buffer); err != nil {
	// 		return nil, err
	// 	}
	// 	var rec = *p_rec
	// 	d.authorities = append(d.authorities, rec)
	// }
	//
	// for i := 0; i < int(d.Header.resource_entries); i++ {
	//
	// 	var p_rec *DnsRecord
	// 	if p_rec, err = eDnsRecord.Read(*buffer); err != nil {
	// 		return nil, err
	// 	}
	// 	var rec = *p_rec
	// 	d.resources = append(d.resources, rec)
	// }

	return d, nil
}
