package core

type ResultCode uint

const (
	NOERROR ResultCode = iota
	FORMERR
	SERVFAIL
	NXDOMAIN
	NOTIMP
	REFUSED
)

func (res ResultCode) From_num(num uint8) ResultCode {
	switch num {
	case 1:
		return FORMERR
	case 2:
		return SERVFAIL
	case 3:
		return NXDOMAIN
	case 4:
		return NOTIMP
	case 5:
		return REFUSED
	default:
		return NOERROR
	}
}
