package core

import (
	"errors"
	"math/rand"
	"net"

	"github.com/orangeseeds/mothDNS/pkg/bpb"
)

type DNSPacket struct {
	Questions   []DNSQuestion `json:"questions"`
	Answers     []DNSRecord   `json:"answers"`
	Authorities []DNSRecord   `json:"authorities"`
	Resources   []DNSRecord   `json:"resource"`
	Header      dnsHeader     `json:"header"`
}

func NewPacket() DNSPacket {
	d := DNSPacket{
		Header:      *NewHeader(),
		Questions:   []DNSQuestion{},
		Answers:     []DNSRecord{},
		Authorities: []DNSRecord{},
		Resources:   []DNSRecord{},
	}

	return d
}

func (d *DNSPacket) From_buffer(buffer *bpb.BytePacketBuffer) (*DNSPacket, error) {
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

		var p_rec *DNSRecord
		if p_rec, err = eDnsRecord.Read(buffer); err != nil {
			return nil, err
		}
		d.Answers = append(d.Answers, *p_rec)
	}
	for i := 0; i < int(d.Header.AuthoritativeEntries); i++ {

		var p_rec *DNSRecord
		if p_rec, err = eDnsRecord.Read(buffer); err != nil {
			return nil, err
		}
		rec := *p_rec
		d.Authorities = append(d.Authorities, rec)
	}
	for i := 0; i < int(d.Header.ResourceEntries); i++ {

		var p_rec *DNSRecord
		if p_rec, err = eDnsRecord.Read(buffer); err != nil {
			return nil, err
		}
		rec := *p_rec
		d.Resources = append(d.Resources, rec)
	}

	return d, nil
}

// ---------------------------------- For Writing ---------------------------------------------------

func (d *DNSPacket) Write(buffer *bpb.BytePacketBuffer) error {
	d.Header.Questions = uint16(len(d.Questions))
	d.Header.Answers = uint16(len(d.Answers))
	d.Header.AuthoritativeEntries = uint16(len(d.Authorities))
	d.Header.ResourceEntries = uint16(len(d.Resources))

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

func (d *DNSPacket) GetRandomA() (net.IP, error) {
	var answers []DNSRecord

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

func (d *DNSPacket) GetResolvedNS(qname string) (string, error) {
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

func (d *DNSPacket) GetUnresNS(qname string) (string, error) {
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
