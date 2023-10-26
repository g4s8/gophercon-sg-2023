package strings

import (
	"unsafe"
)

func bytesToString(b []byte) string {
	return string(b)
}

func bytesToStringUnsafe(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
