package model

import (
	"golang.org/x/text/encoding/simplifiedchinese"
)

// GB18030 实现了dbf.Decoder
type GB18030 struct{}

// Decode 将GB18030转成utf8
func (GB18030) Decode(in []byte) ([]byte, error) {
	dec := simplifiedchinese.GB18030.NewDecoder()
	return dec.Bytes(in)
}
