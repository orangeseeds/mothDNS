package core

import (
	"fmt"
	"net"

	"github.com/orangeseeds/mothDNS/pkg/bpb"
)

type DNSRecord struct {
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
		Host     string `json:"host,omitempty"`
		Priority uint16 `json:"priority,omitempty"`
		Ttl      uint32 `json:"ttl,omitempty"`
	} `json:"MX,omitempty"`
	AAAA struct {
		Domain string `json:"domain,omitempty"`
		Addr   net.IP `json:"addr,omitempty"`
		Ttl    uint32 `json:"ttl,omitempty"`
	} `json:"AAAA,omitempty"`
	Type uint `json:"type"`
}

var eDnsRecord = &DNSRecord{}

func (d *DNSRecord) Read(buffer *bpb.BytePacketBuffer) (*DNSRecord, error) {
	var err error
	domain := ""
	err = buffer.ReadQName(&domain)
	if err != nil {
		return nil, err
	}

	var result uint16
	if result, err = buffer.ReadTwoBytes(); err != nil {
		return nil, err
	}
	qtype_num := result

	qtype := QueryType.From_num(0, qtype_num)
	if _, err = buffer.ReadTwoBytes(); err != nil {
		return nil, err
	}

	var result_32 uint32
	if result_32, err = buffer.ReadFourBytes(); err != nil {
		return nil, err
	}
	ttl := result_32

	if result, err = buffer.ReadTwoBytes(); err != nil {
		return nil, err
	}
	data_len := result

	switch qtype {
	case QT_A:
		if result_32, err = buffer.ReadFourBytes(); err != nil {
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
		if err = buffer.ReadQName(&ns); err != nil {
			return nil, err
		}

		eDnsRecord.NS.Domain = domain
		eDnsRecord.NS.Host = ns
		eDnsRecord.NS.Ttl = ttl

		eDnsRecord.Type = 2
		return eDnsRecord, nil

	case QT_CNAME:
		var cname string
		if err = buffer.ReadQName(&cname); err != nil {
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

		if priority, err = buffer.ReadTwoBytes(); err != nil {
			return nil, err
		}

		if err = buffer.ReadQName(&mx); err != nil {
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

func (d *DNSRecord) Write(buffer *bpb.BytePacketBuffer) (uint, error) {
	start_pos := buffer.Pos()

	switch d.Type {
	// For A
	case 1:
		if err := buffer.WriteQName(&d.A.Domain); err != nil {
			return 0, err
		}
		if err := buffer.WriteTwoBytes(uint16(QT_A)); err != nil {
			return 0, err
		}
		if err := buffer.WriteTwoBytes(1); err != nil {
			return 0, err
		}
		if err := buffer.WriteFourBytes(d.A.Ttl); err != nil {
			return 0, err
		}
		if err := buffer.WriteTwoBytes(4); err != nil {
			return 0, err
		}

		octets := d.A.Addr[len(d.A.Addr)-4:]
		if err := buffer.WriteOneByte(octets[0]); err != nil {
			return 0, err
		}
		if err := buffer.WriteOneByte(octets[1]); err != nil {
			return 0, err
		}
		if err := buffer.WriteOneByte(octets[2]); err != nil {
			return 0, err
		}
		if err := buffer.WriteOneByte(octets[3]); err != nil {
			return 0, err
		}

	case 2:
		if err := buffer.WriteQName(&d.NS.Domain); err != nil {
			return 0, err
		}
		if err := buffer.WriteTwoBytes(uint16(QT_NS)); err != nil {
			return 0, err
		}
		if err := buffer.WriteTwoBytes(1); err != nil {
			return 0, err
		}
		if err := buffer.WriteFourBytes(d.NS.Ttl); err != nil {
			return 0, err
		}

		pos := buffer.Pos()
		if err := buffer.WriteTwoBytes(0); err != nil {
			return 0, err
		}

		if err := buffer.WriteQName(&d.NS.Host); err != nil {
			return 0, err
		}

		size := buffer.Pos() - (pos + 2)
		if err := buffer.SetTwoBytes(pos, uint16(size)); err != nil {
			return 0, err
		}

	case 5:
		if err := buffer.WriteQName(&d.CNAME.Domain); err != nil {
			return 0, err
		}
		if err := buffer.WriteTwoBytes(uint16(QT_CNAME)); err != nil {
			return 0, err
		}
		if err := buffer.WriteTwoBytes(1); err != nil {
			return 0, err
		}
		if err := buffer.WriteFourBytes(d.CNAME.Ttl); err != nil {
			return 0, err
		}

		pos := buffer.Pos()
		if err := buffer.WriteTwoBytes(0); err != nil {
			return 0, err
		}

		if err := buffer.WriteQName(&d.NS.Host); err != nil {
			return 0, err
		}

		size := buffer.Pos() - (pos + 2)
		if err := buffer.SetTwoBytes(pos, uint16(size)); err != nil {
			return 0, err
		}

	case 15:
		if err := buffer.WriteQName(&d.MX.Domain); err != nil {
			return 0, err
		}
		if err := buffer.WriteTwoBytes(uint16(QT_MX)); err != nil {
			return 0, err
		}
		if err := buffer.WriteTwoBytes(1); err != nil {
			return 0, err
		}
		if err := buffer.WriteFourBytes(d.MX.Ttl); err != nil {
			return 0, err
		}

		pos := buffer.Pos()
		if err := buffer.WriteTwoBytes(0); err != nil {
			return 0, err
		}

		if err := buffer.WriteTwoBytes(d.MX.Priority); err != nil {
			return 0, err
		}

		if err := buffer.WriteQName(&d.MX.Host); err != nil {
			return 0, err
		}

		size := buffer.Pos() - (pos + 2)
		if err := buffer.SetTwoBytes(pos, uint16(size)); err != nil {
			return 0, err
		}

	case 28:
		if err := buffer.WriteQName(&d.AAAA.Domain); err != nil {
			return 0, err
		}
		if err := buffer.WriteTwoBytes(uint16(QT_AAAA)); err != nil {
			return 0, err
		}
		if err := buffer.WriteTwoBytes(1); err != nil {
			return 0, err
		}
		if err := buffer.WriteFourBytes(d.AAAA.Ttl); err != nil {
			return 0, err
		}
		if err := buffer.WriteTwoBytes(16); err != nil {
			return 0, err
		}

		for octet := range d.AAAA.Addr {
			if err := buffer.WriteTwoBytes(uint16(octet)); err != nil {
				return 0, err
			}
		}

	// For UNKONWN
	default:
		fmt.Printf("Skipping record %+v\n", d.UNKNOWN)
	}
	return (buffer.Pos() - start_pos), nil
}
