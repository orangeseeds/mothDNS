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
	case qt_A:
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
	case qt_UNKNOWN:
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
	// For UNKONWN
	case 0:
		fmt.Printf("Skipping record %+v", d.UNKNOWN)
		// For A
	case 1:
		if err := buffer.Write_qname(&d.A.domain); err != nil {
			return 0, err
		}
		if err := buffer.Write_u16(uint16(qt_A)); err != nil {
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

	}
	return (buffer.Pos() - start_pos), nil
}
