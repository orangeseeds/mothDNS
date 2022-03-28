package header

import (
	// "errors"
	"fmt"
	// "strings"
	"net"
)

type DnsRecord struct {
	UNKNOWN struct {
		domain   string
		qtype    uint16
		data_len uint16
		ttl      uint32
	}
	A struct {
		domain string
		addr   net.IP
		ttl    uint32
	}
	NS struct {
		domain string
		host   string
		ttl    uint32
	}
	CNAME struct {
		domain string
		host   string
		ttl    uint32
	}
	MX struct {
		domain   string
		priority uint16
		host     string
		ttl      uint32
	}
	AAAA struct {
		domain string
		addr   net.IP
		ttl    uint32
	}
	Type uint
}

var eDnsRecord = &DnsRecord{}

func (d *DnsRecord) Read(buffer *BytePacketBuffer) (*DnsRecord, error) {

	// fmt.Println("posi record", buffer.pos)
	var err error
	domain := ""
	buffer.Read_qname(&domain)

	var result uint16
	if result, err = buffer.Read_u16(); err != nil {
		return nil, err
	}
	qtype_num := result
	// fmt.Println("q_type_num---------", qtype_num)

	qtype := QueryType.From_num(0, qtype_num)
	// fmt.Println("q_type---------", qtype)
	if _, err = buffer.Read_u16(); err != nil {
		return nil, err
	}

	// fmt.Println(domain)

	var result_32 uint32
	if result_32, err = buffer.Read_u32(); err != nil {
		return nil, err
	}
	ttl := result_32
	// fmt.Println("ttl---------", ttl)

	if result, err = buffer.Read_u16(); err != nil {
		return nil, err
	}
	data_len := result
	// fmt.Println("data_len---------", data_len)

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

		eDnsRecord.A.domain = domain
		eDnsRecord.A.addr = addr

		eDnsRecord.A.ttl = ttl
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

		eDnsRecord.NS.domain = domain
		eDnsRecord.NS.host = ns
		eDnsRecord.NS.ttl = ttl

		return eDnsRecord, nil

	case QT_CNAME:
		var cname string
		if err = buffer.Read_qname(&cname); err != nil {
			return nil, err
		}

		eDnsRecord.CNAME.domain = domain
		eDnsRecord.CNAME.host = cname
		eDnsRecord.CNAME.ttl = ttl

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

		eDnsRecord.MX.domain = domain
		eDnsRecord.MX.priority = priority
		eDnsRecord.MX.host = mx
		eDnsRecord.MX.ttl = ttl

		return eDnsRecord, nil

	case QT_UNKNOWN:
		if err = buffer.Step(uint(data_len)); err != nil {
			return nil, err
		}

		eDnsRecord.UNKNOWN.domain = domain
		eDnsRecord.UNKNOWN.qtype = qtype_num
		eDnsRecord.UNKNOWN.data_len = data_len
		eDnsRecord.UNKNOWN.ttl = ttl
		eDnsRecord.Type = 0

		return eDnsRecord, nil
	}
	return nil, nil
}

// ################################## For Writing ###################################################

func (d *DnsRecord) Write(buffer *BytePacketBuffer) (uint, error) {

	start_pos := buffer.Pos()

	switch d.Type {
	// For A
	case 1:
		if err := buffer.Write_qname(&d.A.domain); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(uint16(QT_A)); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(1); err != nil {
			return 0, err
		}
		if err := buffer.Write_u32(d.A.ttl); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(4); err != nil {
			return 0, err
		}

		octets := d.A.addr[len(d.A.addr)-4:]
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
		if err := buffer.Write_qname(&d.NS.domain); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(uint16(QT_NS)); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(1); err != nil {
			return 0, err
		}
		if err := buffer.Write_u32(d.NS.ttl); err != nil {
			return 0, err
		}

		pos := buffer.Pos()
		if err := buffer.Write_u16(0); err != nil {
			return 0, err
		}

		if err := buffer.Write_qname(&d.NS.host); err != nil {
			return 0, err
		}

		size := buffer.Pos() - (pos + 2)
		if err := buffer.Set_u16(pos, uint16(size)); err != nil {
			return 0, err
		}

	case 5:
		if err := buffer.Write_qname(&d.CNAME.domain); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(uint16(QT_CNAME)); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(1); err != nil {
			return 0, err
		}
		if err := buffer.Write_u32(d.CNAME.ttl); err != nil {
			return 0, err
		}

		pos := buffer.Pos()
		if err := buffer.Write_u16(0); err != nil {
			return 0, err
		}

		if err := buffer.Write_qname(&d.NS.host); err != nil {
			return 0, err
		}

		size := buffer.Pos() - (pos + 2)
		if err := buffer.Set_u16(pos, uint16(size)); err != nil {
			return 0, err
		}

	case 15:
		if err := buffer.Write_qname(&d.MX.domain); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(uint16(QT_MX)); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(1); err != nil {
			return 0, err
		}
		if err := buffer.Write_u32(d.MX.ttl); err != nil {
			return 0, err
		}

		pos := buffer.Pos()
		if err := buffer.Write_u16(0); err != nil {
			return 0, err
		}

		if err := buffer.Write_u16(d.MX.priority); err != nil {
			return 0, err
		}

		if err := buffer.Write_qname(&d.MX.host); err != nil {
			return 0, err
		}

		size := buffer.Pos() - (pos + 2)
		if err := buffer.Set_u16(pos, uint16(size)); err != nil {
			return 0, err
		}

	case 28:
		if err := buffer.Write_qname(&d.AAAA.domain); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(uint16(QT_AAAA)); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(1); err != nil {
			return 0, err
		}
		if err := buffer.Write_u32(d.AAAA.ttl); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(16); err != nil {
			return 0, err
		}

		for octet := range d.AAAA.addr {
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
