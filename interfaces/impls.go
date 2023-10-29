package main

import "strconv"

type inter64 int64

type (
	inter32     int32
	interStr    string
	interPtr    int64
	interStruct struct {
		val int64
	}
	interPtrStruct struct {
		val int64
	}
	interStructPtr struct {
		val *int64
	}
	interPtrStructPtr struct {
		val *int64
	}
)

func (i inter64) Int64() int64 {
	return int64(i)
}

func (i inter32) Int64() int64 {
	return int64(i)
}

func (i interStr) Int64() int64 {
	val, _ := strconv.ParseInt(string(i), 10, 64)
	return val
}

func (i *interPtr) Int64() int64 {
	return int64(*i)
}

func (i interStruct) Int64() int64 {
	return i.val
}

func (i *interPtrStruct) Int64() int64 {
	return i.val
}

func (i interStructPtr) Int64() int64 {
	return *i.val
}

func (i *interPtrStructPtr) Int64() int64 {
	return *i.val
}
