package strings

import "unsafe"

type MutableString []byte

func NewMutableString(s string) MutableString {
	return MutableString([]byte(s))
}

func (s MutableString) String() string {
	return unsafe.String(unsafe.SliceData(s), len(s))
}

func (s MutableString) Bytes() []byte {
	return s
}
