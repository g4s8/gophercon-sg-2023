package strings

import (
	"bytes"
	"strings"
	"unsafe"
)

func concatStringsAdd(ss []string) string {
	var s string
	for _, v := range ss {
		s += v
	}
	return s
}

func concatStringsJoin(ss []string) string {
	return strings.Join(ss, "")
}

func concatStringsBuilder(ss []string) string {
	var b strings.Builder
	for _, v := range ss {
		b.WriteString(v)
	}
	return b.String()
}

func concatStringsBuilderGrow(ss []string) string {
	size := allStrBytesLen(ss)
	var b strings.Builder
	b.Grow(size)
	for _, v := range ss {
		b.WriteString(v)
	}
	return b.String()
}

func concatStringsBuffer(ss []string) string {
	var b bytes.Buffer
	for _, v := range ss {
		b.WriteString(v)
	}
	return b.String()
}

func concatStringsBufferGrow(ss []string) string {
	size := allStrBytesLen(ss)
	var b bytes.Buffer
	b.Grow(size)
	for _, v := range ss {
		b.WriteString(v)
	}
	return b.String()
}

func concatStringsCopy(ss []string) string {
	size := allStrBytesLen(ss)
	b := make([]byte, size)
	var i int
	for _, v := range ss {
		i += copy(b[i:], v)
	}
	return string(b)
}

func concatStringsAppend(ss []string) string {
	size := allStrBytesLen(ss)
	b := make([]byte, 0, size)
	for _, v := range ss {
		b = append(b, v...)
	}
	return string(b)
}

func concatStringsCopyUnsafe(ss []string) string {
	size := allStrBytesLen(ss)
	b := make([]byte, size)
	var i int
	for _, v := range ss {
		i += copy(b[i:], v)
	}
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func concatStringsAppendUnsafe(ss []string) string {
	size := allStrBytesLen(ss)
	b := make([]byte, 0, size)
	for _, v := range ss {
		b = append(b, v...)
	}
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func concatStringsCopyUnsafe2(ss []string) string {
	strlen := allStrBytesLen(ss)
	b := make([]byte, strlen)
	var pos int
	for _, v := range ss {
		l := len(v)
		ptr := unsafe.StringData(v)
		bts := unsafe.Slice(ptr, l)
		pos += copy(b[pos:], bts)
	}
	return unsafe.String(unsafe.SliceData(b), strlen)
}

// cached version

// builder, buffer, copy, copyUnsafe, copyUnsafe2

type builderConcater struct {
	builder strings.Builder
}

func (c *builderConcater) concat(ss []string) string {
	c.builder.Reset()
	n := allStrBytesLen(ss)
	c.builder.Grow(n)

	for _, v := range ss {
		c.builder.WriteString(v)
	}
	return c.builder.String()
}

type bufferConcater struct {
	buffer bytes.Buffer
}

func (c *bufferConcater) concat(ss []string) string {
	c.buffer.Reset()
	n := allStrBytesLen(ss)
	c.buffer.Grow(n)

	for _, v := range ss {
		c.buffer.WriteString(v)
	}
	return c.buffer.String()
}

type copyConcater struct {
	b []byte
}

func (c *copyConcater) concat(ss []string) string {
	n := allStrBytesLen(ss)
	if cap(c.b) < n {
		c.b = make([]byte, n)
	} else {
		c.b = c.b[:n]
	}
	var i int
	for _, v := range ss {
		i += copy(c.b[i:], v)
	}
	return string(c.b)
}

type copyUnsafeConcater struct{ b []byte }

func (c *copyUnsafeConcater) concat(ss []string) string {
	n := allStrBytesLen(ss)
	if cap(c.b) < n {
		c.b = make([]byte, n)
	} else {
		c.b = c.b[:n]
	}
	var i int
	for _, v := range ss {
		i += copy(c.b[i:], v)
	}
	return unsafe.String(unsafe.SliceData(c.b), n)
}

type appendUnsafeConcater struct{ b []byte }

func (c *appendUnsafeConcater) concat(ss []string) string {
	n := allStrBytesLen(ss)
	if cap(c.b) < n {
		c.b = make([]byte, 0, n)
	} else {
		c.b = c.b[:0]
	}
	for _, v := range ss {
		c.b = append(c.b, v...)
	}
	return unsafe.String(unsafe.SliceData(c.b), n)
}

type copyUnsafe2Concater struct {
	b []byte
}

func (c *copyUnsafe2Concater) concat(ss []string) string {
	n := allStrBytesLen(ss)
	if cap(c.b) < n {
		c.b = make([]byte, n)
	} else {
		c.b = c.b[:n]
	}
	var pos int
	for _, v := range ss {
		l := len(v)
		ptr := unsafe.StringData(v)
		bts := unsafe.Slice(ptr, l)
		pos += copy(c.b[pos:], bts)
	}
	return unsafe.String(unsafe.SliceData(c.b), n)
}

func allStrBytesLen(ss []string) int {
	var n int
	for _, v := range ss {
		n += len(v)
	}
	return n
}
