package core

import (
	"errors"
	"math/rand"
	"net"
)

// "fmt"

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
	// answers := map[string]string{}

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

// func (d *DnsPacket) GetResolvedNS(qname string) (net.IP, error) {
// 	nsMap, _ := d.GetNS(qname)
// 	for _, ns := range nsMap {
// 		var aResources []DnsRecord
// 		for _, r := range d.Resources {
// 			if r.A.Domain != "" {
// 				aResources = append(aResources, r)
// 			}
// 		}
// 		for _, a := range aResources {
// 			if a.A.Domain == nsMap[ns] {
// 				return a.A.Addr, nil
// 			}
// 		}

// 	}
// 	return nil, errors.New("no ip found for " + qname)
// }

// func (d *DnsPacket) GetUnresolvedNS(qname string) (string, error) {

// 	nsMap, _ := d.GetNS(qname)
// 	for _, ns := range nsMap {
// 		return nsMap[ns], nil
// 	}
// 	return "", errors.New(qname + " not in nsMap")
// }
