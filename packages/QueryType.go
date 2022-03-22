package header

import (
// "errors"
// "fmt"
// "strings"
)

type QueryType uint16

const (
	qt_UNKNOWN QueryType = iota
	qt_A
)

func (q QueryType) To_num() uint16 {
	switch q {
	case qt_A:
		return 1
	}
	return uint16(qt_UNKNOWN)
}

func (q QueryType) From_num(num uint16) QueryType {
	switch num {
	case 1:
		return qt_A
	}
	return qt_UNKNOWN
}

// func main() {
// 	println(qt_UNKNOWN)
// }
