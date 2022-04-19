package core

import (
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
type SOA struct {
	Domain  string
	Mname   string
	Rname   string
	Serial  uint32
	Refresh uint32
	Retry   uint32
	Expire  uint32
	Minimum uint32
	Ttl     uint32
}
type TXT struct {
	Domain string
	Data   string
	Ttl    uint32
}

func (r UNKNOWN) recordType() {}
func (r A) recordType()       {}
func (r NS) recordType()      {}
func (r CNAME) recordType()   {}
func (r MX) recordType()      {}
func (r AAAA) recordType()    {}
func (r SOA) recordType()     {}

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
	case QT_SOA:
		var serial, refresh, retry, expire, minimum uint32
		mName := ""
		if err = buffer.Read_qname(&mName); err != nil {
			return nil, err
		}
		rName := ""
		if err = buffer.Read_qname(&rName); err != nil {
			return nil, err
		}
		if serial, err = buffer.Read_u32(); err != nil {
			return nil, err
		}
		if refresh, err = buffer.Read_u32(); err != nil {
			return nil, err
		}
		if retry, err = buffer.Read_u32(); err != nil {
			return nil, err
		}
		if expire, err = buffer.Read_u32(); err != nil {
			return nil, err
		}
		if minimum, err = buffer.Read_u32(); err != nil {
			return nil, err
		}

		s := SOA{
			Domain:  domain,
			Mname:   mName,
			Rname:   rName,
			Serial:  serial,
			Refresh: refresh,
			Retry:   retry,
			Expire:  expire,
			Minimum: minimum,
			Ttl:     ttl,
		}
		return s, nil

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
	case SOA:
		if err := buffer.Write_qname(r.(SOA).Domain); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(1); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(uint16(r.(SOA).Ttl)); err != nil {
			return 0, err
		}

		pos := buffer.Pos()
		if err := buffer.Write_u16(0); err != nil {
			return 0, err
		}

		if err := buffer.Write_qname(r.(SOA).Mname); err != nil {
			return 0, err
		}

		if err := buffer.Write_qname(r.(SOA).Rname); err != nil {
			return 0, err
		}
		if err := buffer.Write_u32(r.(SOA).Serial); err != nil {
			return 0, err
		}
		if err := buffer.Write_u32(r.(SOA).Refresh); err != nil {
			return 0, err
		}
		if err := buffer.Write_u32(r.(SOA).Retry); err != nil {
			return 0, err
		}

		if err := buffer.Write_u32(r.(SOA).Expire); err != nil {
			return 0, err
		}
		if err := buffer.Write_u32(r.(SOA).Minimum); err != nil {
			return 0, err
		}

		size := buffer.Pos() - (pos + 2)
		buffer.Set_u16(pos, uint16(size))
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
		// log.Printf("Skipping record of query type %d\n", r.(UNKNOWN).Qtype)
		// fmt.Printf("Skipping record %+v", r.(UNKNOWN))
	}
	return (buffer.Pos() - start_pos), nil
}

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

// case 28:
// 	if err := buffer.Write_qname(d.AAAA.Domain); err != nil {
// 		return 0, err
// 	}
// 	if err := buffer.Write_u16(uint16(QT_AAAA)); err != nil {
// 		return 0, err
// 	}
// 	if err := buffer.Write_u16(1); err != nil {
// 		return 0, err
// 	}
// 	if err := buffer.Write_u32(d.AAAA.Ttl); err != nil {
// 		return 0, err
// 	}
// 	if err := buffer.Write_u16(16); err != nil {
// 		return 0, err
// 	}

// 	for octet := range d.AAAA.Addr {
// 		if err := buffer.Write_u16(uint16(octet)); err != nil {
// 			return 0, err
// 		}
// 	}
