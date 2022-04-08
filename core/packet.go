package core

import (
// "fmt"
)

type DnsPacket struct {
	Header      DnsHeader     `json:"header"`
	Questions   []DnsQuestion `json:"questions"`
	Answers     []DnsRecord   `json:"answers"`
	Authorities []DnsRecord   `json:"authorities"`
	Resources   []DnsRecord   `json:"resource"`
}

func NewPacket() DnsPacket {
	d := DnsPacket{
		Header:      *NewHeader(),
		Questions:   []DnsQuestion{},
		Answers:     []DnsRecord{},
		Authorities: []DnsRecord{},
		Resources:   []DnsRecord{},
	}

	return d
}

func (d *DnsPacket) From_buffer(buffer *BytePacketBuffer) (*DnsPacket, error) {
	// fmt.Println(buffer)
	var err error

	if err = d.Header.Read(buffer); err != nil {
		return nil, err
	}
	// fmt.Println(d.Header)

	// fmt.Println("DNS Questions")
	for i := 0; i < int(d.Header.Questions); i++ {

		question := NewQuestion("", QT_UNKNOWN)
		if err = question.Read(buffer); err != nil {
			return nil, err
		}

		d.Questions = append(d.Questions, question)
		// fmt.Println(buffer.pos)
	}
	// fmt.Println(d.Questions)
	//
	// fmt.Println("DNS Answers")
	for i := 0; i < int(d.Header.Answers); i++ {

		var p_rec *DnsRecord
		if p_rec, err = eDnsRecord.Read(buffer); err != nil {
			return nil, err
		}
		d.Answers = append(d.Answers, *p_rec)
	}
	// fmt.Println(d.Answers)
	//
	for i := 0; i < int(d.Header.Authoritative_entries); i++ {

		var p_rec *DnsRecord
		if p_rec, err = eDnsRecord.Read(buffer); err != nil {
			return nil, err
		}
		var rec = *p_rec
		d.Authorities = append(d.Authorities, rec)
	}
	//
	for i := 0; i < int(d.Header.Resource_entries); i++ {

		var p_rec *DnsRecord
		if p_rec, err = eDnsRecord.Read(buffer); err != nil {
			return nil, err
		}
		var rec = *p_rec
		d.Resources = append(d.Resources, rec)
	}

	return d, nil
}

// ################################## For Writing ###################################################

func (d *DnsPacket) Write(buffer *BytePacketBuffer) error {

	// if err := buffer.Write_u16(d.id); err != nil {
	// 	return err
	// }

	d.Header.Questions = uint16(len(d.Questions))
	d.Header.Answers = uint16(len(d.Answers))
	d.Header.Authoritative_entries = uint16(len(d.Authorities))
	d.Header.Resource_entries = uint16(len(d.Resources))

	if err := d.Header.Write(buffer); err != nil {
		return err
	}

	for _, question := range d.Questions {
		if err := question.Write(buffer); err != nil {
			return err
		}
	}

	for _, rec := range d.Answers {
		if _, err := rec.Write(buffer); err != nil {
			return err
		}
	}

	for _, rec := range d.Authorities {
		if _, err := rec.Write(buffer); err != nil {
			return err
		}
	}

	for _, rec := range d.Resources {
		if _, err := rec.Write(buffer); err != nil {
			return err
		}
	}

	return nil
}
