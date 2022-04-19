package core

import (
	"fmt"
	"net"
)

type Record interface {
	recordType()
}

type UNKNOWN struct {
	Domain   string `json:"domain,omitempty"`
	Qtype    uint16 `json:"qtype,omitempty"`
	Data_len uint16 `json:"data_len,omitempty"`
	Ttl      uint32 `json:"ttl,omitempty"`
}
type A struct {
	Domain string `json:"domain,omitempty"`
	Addr   net.IP `json:"addr,omitempty"`
	Ttl    uint32 `json:"ttl,omitempty"`
}
type NS struct {
	Domain string `json:"domain,omitempty"`
	Host   string `json:"host,omitempty"`
	Ttl    uint32 `json:"ttl,omitempty"`
}
type CNAME struct {
	Domain string `json:"domain,omitempty"`
	Host   string `json:"host,omitempty"`
	Ttl    uint32 `json:"ttl,omitempty"`
}
type MX struct {
	Domain   string `json:"domain,omitempty"`
	Priority uint16 `json:"priority,omitempty"`
	Host     string `json:"host,omitempty"`
	Ttl      uint32 `json:"ttl,omitempty"`
}
type AAAA struct {
	Domain string `json:"domain,omitempty"`
	Addr   net.IP `json:"addr,omitempty"`
	Ttl    uint32 `json:"ttl,omitempty"`
}

func (r UNKNOWN) recordType() {}
func (r A) recordType()       {}
func (r NS) recordType()      {}
func (r CNAME) recordType()   {}
func (r MX) recordType()      {}
func (r AAAA) recordType()    {}

func ReadRecord(buffer *BytePacketBuffer) (Record, error) {

	var err error
	domain := ""
	buffer.Read_qname(&domain)

	var result uint16
	if result, err = buffer.Read_u16(); err != nil {
		return nil, err
	}
	qtype_num := result

	qtype := QueryType.From_num(0, qtype_num)
	if _, err = buffer.Read_u16(); err != nil {
		return nil, err
	}

	var result_32 uint32
	if result_32, err = buffer.Read_u32(); err != nil {
		return nil, err
	}
	ttl := result_32

	if result, err = buffer.Read_u16(); err != nil {
		return nil, err
	}
	data_len := result

	switch qtype {
	case QT_A:
		if result_32, err = buffer.Read_u32(); err != nil {
			return nil, err
		}
		raw_addr := result_32
		p1 := uint8((raw_addr >> 24) & 0xFF)
		p2 := uint8((raw_addr >> 16) & 0xFF)
		p3 := uint8((raw_addr >> 8) & 0xFF)
		p4 := uint8((raw_addr >> 0) & 0xFF)
		addr := net.IPv4(p1, p2, p3, p4)
		a := A{
			Domain: domain,
			Addr:   addr,
			Ttl:    ttl,
		}
		return a, nil
	case QT_NS:
		var ns string
		if err = buffer.Read_qname(&ns); err != nil {
			return nil, err
		}
		n := NS{
			Domain: domain,
			Host:   ns,
			Ttl:    ttl,
		}
		return n, nil
	case QT_CNAME:
		var cname string
		if err = buffer.Read_qname(&cname); err != nil {
			return nil, err
		}
		c := CNAME{
			Domain: domain,
			Host:   cname,
			Ttl:    ttl,
		}
		return c, nil
	case QT_MX:
		var priority uint16
		var mx string
		if priority, err = buffer.Read_u16(); err != nil {
			return nil, err
		}
		if err = buffer.Read_qname(&mx); err != nil {
			return nil, err
		}
		m := MX{
			Domain:   domain,
			Priority: priority,
			Host:     mx,
			Ttl:      ttl,
		}
		return m, nil
	default:
		if err = buffer.Step(uint(data_len)); err != nil {
			return nil, err
		}
		u := UNKNOWN{
			Domain:   domain,
			Qtype:    qtype_num,
			Data_len: data_len,
			Ttl:      ttl,
		}
		return u, nil
	}
}

func WriteRecord(r Record, buffer *BytePacketBuffer) (uint, error) {

	start_pos := buffer.Pos()

	switch r.(type) {
	// For A
	case A:
		if err := buffer.Write_qname(r.(A).Domain); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(uint16(QT_A)); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(1); err != nil {
			return 0, err
		}
		if err := buffer.Write_u32(r.(A).Ttl); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(4); err != nil {
			return 0, err
		}

		octets := r.(A).Addr[len(r.(A).Addr)-4:]
		if err := buffer.Write_u8(octets[0]); err != nil {
			return 0, err
		}
		if err := buffer.Write_u8(octets[1]); err != nil {
			return 0, err
		}
		if err := buffer.Write_u8(octets[2]); err != nil {
			return 0, err
		}
		if err := buffer.Write_u8(octets[3]); err != nil {
			return 0, err
		}

	case NS:
		if err := buffer.Write_qname(r.(NS).Domain); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(uint16(QT_NS)); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(1); err != nil {
			return 0, err
		}
		if err := buffer.Write_u32(r.(NS).Ttl); err != nil {
			return 0, err
		}

		pos := buffer.Pos()
		if err := buffer.Write_u16(0); err != nil {
			return 0, err
		}

		if err := buffer.Write_qname(r.(NS).Host); err != nil {
			return 0, err
		}

		size := buffer.Pos() - (pos + 2)
		if err := buffer.Set_u16(pos, uint16(size)); err != nil {
			return 0, err
		}

	case CNAME:
		if err := buffer.Write_qname(r.(CNAME).Domain); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(uint16(QT_CNAME)); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(1); err != nil {
			return 0, err
		}
		if err := buffer.Write_u32(r.(CNAME).Ttl); err != nil {
			return 0, err
		}

		pos := buffer.Pos()
		if err := buffer.Write_u16(0); err != nil {
			return 0, err
		}

		if err := buffer.Write_qname(r.(NS).Host); err != nil {
			return 0, err
		}

		size := buffer.Pos() - (pos + 2)
		if err := buffer.Set_u16(pos, uint16(size)); err != nil {
			return 0, err
		}

	case MX:
		if err := buffer.Write_qname(r.(MX).Domain); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(uint16(QT_MX)); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(1); err != nil {
			return 0, err
		}
		if err := buffer.Write_u32(r.(MX).Ttl); err != nil {
			return 0, err
		}

		pos := buffer.Pos()
		if err := buffer.Write_u16(0); err != nil {
			return 0, err
		}

		if err := buffer.Write_u16(r.(MX).Priority); err != nil {
			return 0, err
		}

		if err := buffer.Write_qname(r.(MX).Host); err != nil {
			return 0, err
		}

		size := buffer.Pos() - (pos + 2)
		if err := buffer.Set_u16(pos, uint16(size)); err != nil {
			return 0, err
		}

	case AAAA:
		if err := buffer.Write_qname(r.(AAAA).Domain); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(uint16(QT_AAAA)); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(1); err != nil {
			return 0, err
		}
		if err := buffer.Write_u32(r.(AAAA).Ttl); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(16); err != nil {
			return 0, err
		}

		for octet := range r.(AAAA).Addr {
			if err := buffer.Write_u16(uint16(octet)); err != nil {
				return 0, err
			}
		}

	// For UNKONWN
	default:
		fmt.Printf("Skipping record %+v", r.(UNKNOWN))
	}
	return (buffer.Pos() - start_pos), nil
}

type DnsRecord struct {
	UNKNOWN struct {
		Domain   string `json:"domain,omitempty"`
		Qtype    uint16 `json:"qtype,omitempty"`
		Data_len uint16 `json:"data_len,omitempty"`
		Ttl      uint32 `json:"ttl,omitempty"`
	} `json:"UNKNOWN,omitempty"`
	A struct {
		Domain string `json:"domain,omitempty"`
		Addr   net.IP `json:"addr,omitempty"`
		Ttl    uint32 `json:"ttl,omitempty"`
	} `json:"A,omitempty"`
	NS struct {
		Domain string `json:"domain,omitempty"`
		Host   string `json:"host,omitempty"`
		Ttl    uint32 `json:"ttl,omitempty"`
	} `json:"NS,omitempty"`
	CNAME struct {
		Domain string `json:"domain,omitempty"`
		Host   string `json:"host,omitempty"`
		Ttl    uint32 `json:"ttl,omitempty"`
	} `json:"CNAME,omitempty"`
	MX struct {
		Domain   string `json:"domain,omitempty"`
		Priority uint16 `json:"priority,omitempty"`
		Host     string `json:"host,omitempty"`
		Ttl      uint32 `json:"ttl,omitempty"`
	} `json:"MX,omitempty"`
	AAAA struct {
		Domain string `json:"domain,omitempty"`
		Addr   net.IP `json:"addr,omitempty"`
		Ttl    uint32 `json:"ttl,omitempty"`
	} `json:"AAAA,omitempty"`
	Type uint `json:"type"`
}

var eDnsRecord = &DnsRecord{}

func (d *DnsRecord) Read(buffer *BytePacketBuffer) (*DnsRecord, error) {

	var err error
	domain := ""
	buffer.Read_qname(&domain)

	var result uint16
	if result, err = buffer.Read_u16(); err != nil {
		return nil, err
	}
	qtype_num := result

	qtype := QueryType.From_num(0, qtype_num)
	if _, err = buffer.Read_u16(); err != nil {
		return nil, err
	}

	var result_32 uint32
	if result_32, err = buffer.Read_u32(); err != nil {
		return nil, err
	}
	ttl := result_32

	if result, err = buffer.Read_u16(); err != nil {
		return nil, err
	}
	data_len := result

	switch qtype {
	case QT_A:
		if result_32, err = buffer.Read_u32(); err != nil {
			return nil, err
		}
		raw_addr := result_32

		p1 := uint8((raw_addr >> 24) & 0xFF)
		p2 := uint8((raw_addr >> 16) & 0xFF)
		p3 := uint8((raw_addr >> 8) & 0xFF)
		p4 := uint8((raw_addr >> 0) & 0xFF)

		addr := net.IPv4(p1, p2, p3, p4)

		eDnsRecord.A.Domain = domain
		eDnsRecord.A.Addr = addr

		eDnsRecord.A.Ttl = ttl
		eDnsRecord.Type = 1

		return eDnsRecord, nil
	// case QT_AAAA:
	// 	var raw_addr1 uint32
	// 	var raw_addr2 uint32
	// 	var raw_addr3 uint32
	// 	var raw_addr4 uint32
	//
	// 	if raw_addr1, err = buffer.Read_u32(); err != nil {
	// 		return nil, err
	// 	}
	// 	if raw_addr2, err = buffer.Read_u32(); err != nil {
	// 		return nil, err
	// 	}
	// 	if raw_addr3, err = buffer.Read_u32(); err != nil {
	// 		return nil, err
	// 	}
	// 	if raw_addr4, err = buffer.Read_u32(); err != nil {
	// 		return nil, err
	// 	}

	/*
		uint16((raw_addr1 >> 16) & 0xFFFF)
		uint16((raw_addr1 >> 0) & 0xFFFF)
		uint16((raw_addr2 >> 16) & 0xFFFF)
		uint16((raw_addr2 >> 0) & 0xFFFF)
		uint16((raw_addr3 >> 16) & 0xFFFF)
		uint16((raw_addr3 >> 0) & 0xFFFF)
		uint16((raw_addr4 >> 16) & 0xFFFF)
		uint16((raw_addr4 >> 0) & 0xFFFF)
	*/
	case QT_NS:
		var ns string
		if err = buffer.Read_qname(&ns); err != nil {
			return nil, err
		}

		eDnsRecord.NS.Domain = domain
		eDnsRecord.NS.Host = ns
		eDnsRecord.NS.Ttl = ttl

		eDnsRecord.Type = 2
		return eDnsRecord, nil

	case QT_CNAME:
		var cname string
		if err = buffer.Read_qname(&cname); err != nil {
			return nil, err
		}

		eDnsRecord.CNAME.Domain = domain
		eDnsRecord.CNAME.Host = cname
		eDnsRecord.CNAME.Ttl = ttl
		eDnsRecord.Type = 5

		return eDnsRecord, nil

	case QT_MX:
		var priority uint16
		var mx string

		if priority, err = buffer.Read_u16(); err != nil {
			return nil, err
		}

		if err = buffer.Read_qname(&mx); err != nil {
			return nil, err
		}

		eDnsRecord.MX.Domain = domain
		eDnsRecord.MX.Priority = priority
		eDnsRecord.MX.Host = mx
		eDnsRecord.MX.Ttl = ttl

		eDnsRecord.Type = 15
		return eDnsRecord, nil

	default:
		if err = buffer.Step(uint(data_len)); err != nil {
			return nil, err
		}

		eDnsRecord.UNKNOWN.Domain = domain
		eDnsRecord.UNKNOWN.Qtype = qtype_num
		eDnsRecord.UNKNOWN.Data_len = data_len
		eDnsRecord.UNKNOWN.Ttl = ttl
		eDnsRecord.Type = 0

		return eDnsRecord, nil
	}
}

// ---------------------------------- For Writing ---------------------------------------------------

func (d *DnsRecord) Write(buffer *BytePacketBuffer) (uint, error) {

	start_pos := buffer.Pos()

	switch d.Type {
	// For A
	case 1:
		if err := buffer.Write_qname(d.A.Domain); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(uint16(QT_A)); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(1); err != nil {
			return 0, err
		}
		if err := buffer.Write_u32(d.A.Ttl); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(4); err != nil {
			return 0, err
		}

		octets := d.A.Addr[len(d.A.Addr)-4:]
		if err := buffer.Write_u8(octets[0]); err != nil {
			return 0, err
		}
		if err := buffer.Write_u8(octets[1]); err != nil {
			return 0, err
		}
		if err := buffer.Write_u8(octets[2]); err != nil {
			return 0, err
		}
		if err := buffer.Write_u8(octets[3]); err != nil {
			return 0, err
		}

	case 2:
		if err := buffer.Write_qname(d.NS.Domain); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(uint16(QT_NS)); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(1); err != nil {
			return 0, err
		}
		if err := buffer.Write_u32(d.NS.Ttl); err != nil {
			return 0, err
		}

		pos := buffer.Pos()
		if err := buffer.Write_u16(0); err != nil {
			return 0, err
		}

		if err := buffer.Write_qname(d.NS.Host); err != nil {
			return 0, err
		}

		size := buffer.Pos() - (pos + 2)
		if err := buffer.Set_u16(pos, uint16(size)); err != nil {
			return 0, err
		}

	case 5:
		if err := buffer.Write_qname(d.CNAME.Domain); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(uint16(QT_CNAME)); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(1); err != nil {
			return 0, err
		}
		if err := buffer.Write_u32(d.CNAME.Ttl); err != nil {
			return 0, err
		}

		pos := buffer.Pos()
		if err := buffer.Write_u16(0); err != nil {
			return 0, err
		}

		if err := buffer.Write_qname(d.NS.Host); err != nil {
			return 0, err
		}

		size := buffer.Pos() - (pos + 2)
		if err := buffer.Set_u16(pos, uint16(size)); err != nil {
			return 0, err
		}

	case 15:
		if err := buffer.Write_qname(d.MX.Domain); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(uint16(QT_MX)); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(1); err != nil {
			return 0, err
		}
		if err := buffer.Write_u32(d.MX.Ttl); err != nil {
			return 0, err
		}

		pos := buffer.Pos()
		if err := buffer.Write_u16(0); err != nil {
			return 0, err
		}

		if err := buffer.Write_u16(d.MX.Priority); err != nil {
			return 0, err
		}

		if err := buffer.Write_qname(d.MX.Host); err != nil {
			return 0, err
		}

		size := buffer.Pos() - (pos + 2)
		if err := buffer.Set_u16(pos, uint16(size)); err != nil {
			return 0, err
		}

	case 28:
		if err := buffer.Write_qname(d.AAAA.Domain); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(uint16(QT_AAAA)); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(1); err != nil {
			return 0, err
		}
		if err := buffer.Write_u32(d.AAAA.Ttl); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(16); err != nil {
			return 0, err
		}

		for octet := range d.AAAA.Addr {
			if err := buffer.Write_u16(uint16(octet)); err != nil {
				return 0, err
			}
		}

	// For UNKONWN
	default:
		fmt.Printf("Skipping record %+v", d.UNKNOWN)
	}
	return (buffer.Pos() - start_pos), nil
}
