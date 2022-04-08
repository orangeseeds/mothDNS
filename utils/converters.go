package utils

import ()

func To_512buffer(buf []byte) [512]byte {
	reply_buff := [512]byte{0}
	for i := range buf {
		reply_buff[i] = buf[i]
	}

	return reply_buff
}

func From_512buffer(buf [512]byte) []byte {
	reply_buff := make([]byte, 512)
	for i := range buf {
		reply_buff[i] = buf[i]
	}

	return reply_buff
}
