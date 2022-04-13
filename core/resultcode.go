package core

import (
	// "errors"
	"fmt"
	// "strings"
)

type ResultCode uint

const (
	NOERROR  = iota
	FORMERR  = 1
	SERVFAIL = 2
	NXDOMAIN = 3
	NOTIMP   = 4
	REFUSED  = 5
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

func main() {
	var result ResultCode
	fmt.Println(result.From_num(1))
}
