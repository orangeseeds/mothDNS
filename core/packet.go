package core

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
)

type Packet struct {
	Header      DnsHeader     `json:"header"`
	Questions   []DnsQuestion `json:"questions"`
	Answers     []Record      `json:"answers"`
	Authorities []Record      `json:"authorities"`
	Resources   []Record      `json:"resource"`
}

func PacketFrmBuff(buffer *BytePacketBuffer) (*Packet, error) {

	var err error

	d := &Packet{}

	if err = d.Header.Read(buffer); err != nil {
		return nil, err
	}

	for i := 0; i < int(d.Header.Questions); i++ {

		question := NewQuestion("", QT_UNKNOWN)
		if err = question.Read(buffer); err != nil {
			return d, err
		}

		d.Questions = append(d.Questions, question)
	}
	for i := 0; i < int(d.Header.Answers); i++ {

		var rec Record
		if rec, err = ReadRecord(buffer); err != nil {
			return d, err
		}
		d.Answers = append(d.Answers, rec)
	}
	for i := 0; i < int(d.Header.Authoritative_entries); i++ {
		var rec Record
		if rec, err = ReadRecord(buffer); err != nil {
			return d, err
		}
		d.Authorities = append(d.Authorities, rec)
	}
	for i := 0; i < int(d.Header.Resource_entries); i++ {

		var rec Record
		if rec, err = ReadRecord(buffer); err != nil {
			return d, err
		}

		d.Resources = append(d.Resources, rec)
	}

	return d, nil
}

func BuffFrmPacket(d Packet) (BytePacketBuffer, error) {

	d.Header.Questions = uint16(len(d.Questions))
	d.Header.Answers = uint16(len(d.Answers))
	d.Header.Authoritative_entries = uint16(len(d.Authorities))
	d.Header.Resource_entries = uint16(len(d.Resources))

	buffer := NewBuffer()

	if err := d.Header.Write(&buffer); err != nil {
		return buffer, err
	}

	for _, question := range d.Questions {
		if err := question.Write(&buffer); err != nil {
			return buffer, err
		}
	}

	for _, rec := range d.Answers {
		if _, err := WriteRecord(rec, &buffer); err != nil {
			return buffer, err
		}
	}

	for _, rec := range d.Authorities {
		if _, err := WriteRecord(rec, &buffer); err != nil {
			return buffer, err
		}
	}

	for _, rec := range d.Resources {
		if _, err := WriteRecord(rec, &buffer); err != nil {
			return buffer, err
		}
	}

	return buffer, nil
}

func (d *Packet) GetRandomA() (Record, error) {
	var answers []Record
	for _, r := range d.Answers {
		if a, ok := r.(A); ok {
			answers = append(answers, a)
		}
	}

	if len(answers) == 0 {
		return nil, fmt.Errorf("no answers in packet to pick randomly")
	}

	return answers[rand.Intn(len(answers))], nil
}

func (d *Packet) GetResolvedNS() (string, error) {
	for _, r := range d.Authorities {
		if a, ok := r.(NS); ok {
			for _, u := range d.Resources {
				if x, ok := u.(A); ok {
					if a.Host == x.Domain {
						return x.Addr.String(), nil
					}
				}
			}
		}
	}
	return "", fmt.Errorf("no as of type NS")
}

func (d *Packet) GetUnresNS() (string, error) {

	auth := map[string]string{}
	for _, a := range d.Authorities {
		if v, ok := a.(NS); ok {
			auth[v.Domain] = v.Host
		}
	}
	return "", nil
}

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
	var err error

	if err = d.Header.Read(buffer); err != nil {
		return nil, err
	}

	for i := 0; i < int(d.Header.Questions); i++ {

		question := NewQuestion("", QT_UNKNOWN)
		if err = question.Read(buffer); err != nil {
			return nil, err
		}

		d.Questions = append(d.Questions, question)
	}
	for i := 0; i < int(d.Header.Answers); i++ {

		var p_rec *DnsRecord
		if p_rec, err = eDnsRecord.Read(buffer); err != nil {
			return nil, err
		}
		d.Answers = append(d.Answers, *p_rec)
	}
	for i := 0; i < int(d.Header.Authoritative_entries); i++ {

		var p_rec *DnsRecord
		if p_rec, err = eDnsRecord.Read(buffer); err != nil {
			return nil, err
		}
		var rec = *p_rec
		d.Authorities = append(d.Authorities, rec)
	}
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

// ---------------------------------- For Writing ---------------------------------------------------

func (d *DnsPacket) Write(buffer *BytePacketBuffer) error {

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

func (d *DnsPacket) GetRandomA() (net.IP, error) {

	var answers []DnsRecord

	for _, a := range d.Answers {
		if a.A.Domain != "" {
			answers = append(answers, a)
		}
	}
	if len(answers) == 0 {
		return nil, errors.New("no as of type A")
	}
	chosenA := answers[rand.Intn(len(answers))]
	return chosenA.A.Addr, nil
}

func (d *DnsPacket) GetResolvedNS(qname string) (string, error) {

	for _, a := range d.Authorities {
		if a.NS.Host != "" {
			for _, r := range d.Resources {
				if a.NS.Host == r.A.Domain {
					return r.A.Addr.String(), nil
				}
			}
		}
	}
	return "", errors.New("no as of type NS")
}

func (d *DnsPacket) GetUnresNS(qname string) (string, error) {

	authorities := map[string]string{}
	for _, a := range d.Authorities {
		if a.NS.Host != "" {
			authorities[a.NS.Domain] = a.NS.Host
		}
	}

	for _, r := range d.Resources {
		if _, ok := authorities[r.A.Domain]; !ok {
			return r.A.Addr.String(), nil
		}

	}
	return "", errors.New("no as of type NS")
}
