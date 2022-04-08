package core

import (
	"errors"
)

var EOBError = errors.New("Buffer greater than or equals 512")
