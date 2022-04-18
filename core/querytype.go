package core

type QueryType uint16

const (
	QT_UNKNOWN QueryType = iota
	QT_A       QueryType = 1
	QT_NS      QueryType = 2
	QT_CNAME   QueryType = 5
	QT_SOA     QueryType = 6
	QT_PTR     QueryType = 12
	QT_MX      QueryType = 15
	QT_TXT     QueryType = 16
	QT_AAAA    QueryType = 28
)

var MapQt = map[QueryType]string{
	1:  "A",
	2:  "NS",
	5:  "CNAME",
	15: "MX",
	28: "AAAA",
	6:  "SOA",
	12: "PTR",
	16: "TXT",
}

func QtName(qType QueryType) string {
	if val, ok := MapQt[qType]; ok {
		return val
	}
	return "UNKNOWN"
}

func (q QueryType) To_num() uint16 {
	switch q {
	case QT_A:
		return 1
	case QT_NS:
		return 2
	case QT_CNAME:
		return 5
	case QT_SOA:
		return 6
	case QT_PTR:
		return 12
	case QT_MX:
		return 15
	case QT_TXT:
		return 16
	case QT_AAAA:
		return 28
	}
	return uint16(QT_UNKNOWN)
}

func (q QueryType) From_num(num uint16) QueryType {
	switch num {
	case 1:
		return QT_A
	case 2:
		return QT_NS
	case 5:
		return QT_CNAME
	case 6:
		return QT_SOA
	case 12:
		return QT_PTR
	case 15:
		return QT_MX
	case 16:
		return QT_TXT
	case 28:
		return QT_AAAA
	}
	return QT_UNKNOWN
}
