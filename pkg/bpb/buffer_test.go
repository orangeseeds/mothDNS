package bpb

import "testing"

func TestReadQName(t *testing.T) {
	var (
		// 		buff     = []byte{134, 42, 129, 128, 0, 1, 0, 1, 0, 0, 0, 0, 6, 103, 111, 111, 103, 108, 101, 3, 99, 111, 109, 0, 0, 1, 0, 1, 192, 12, 0, 1, 0, 1, 0, 0, 1, 37, 0, 4, 216, 58, 211, 142}
		buff     = []byte{21, 19, 1, 32, 0, 1, 0, 0, 0, 0, 0, 1, 6, 103, 111, 111, 103, 108, 101, 3, 99, 111, 109, 0, 0, 1, 0, 1, 0, 0, 41, 4, 208, 0, 0, 0, 0, 0, 12, 0, 10, 0, 8, 196, 137, 125, 115, 27, 20, 201, 15}
		startPos = 12
		qName    string
	)
	b := New()
	b.Buf = buff
	b.Seek(uint(startPos))

	err := b.ReadQName(&qName)
	if err != nil {
		t.Errorf("%q", err)
	} else if qName != "google.com" {
		t.Errorf("qName for the query expecting to be google.com but got %s", qName)
	}
}