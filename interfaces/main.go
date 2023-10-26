package main

import (
	"strconv"
)

type Inter interface {
	Int64() int64
}

type inter64 int64

func (i inter64) Int64() int64 {
	return int64(i)
}

type inter32 int32

func (i inter32) Int64() int64 {
	return int64(i)
}

type interStr string

func (i interStr) Int64() int64 {
	val, _ := strconv.ParseInt(string(i), 10, 64)
	return val
}

type interPtr int64

func (i *interPtr) Int64() int64 {
	return int64(*i)
}

func main() {
	i32_1 := inter32(1)
	_ = toInt(i32_1)
	_ = toInt32(i32_1)
	_ = toIntGeneric(i32_1)

	i32_256 := inter32(256)
	_ = toInt(i32_256)
	_ = toInt32(i32_256)
	_ = toIntGeneric(i32_256)

	i64_1 := inter64(1)
	_ = toInt(i64_1)
	_ = toInt64(i64_1)
	_ = toIntGeneric(i64_1)

	i64_256 := inter64(256)
	_ = toInt(i64_256)
	_ = toInt64(i64_256)
	_ = toIntGeneric(i64_256)

	iStr := interStr("1")
	_ = toInt(iStr)
	_ = toIntStr(iStr)
	_ = toIntGeneric(iStr)

	iPtr := interPtr(1)
	_ = toInt(&iPtr)
	_ = toIntPtr(&iPtr)
	_ = toIntGeneric(&iPtr)
}

func toInt(i Inter) int64 {
	return i.Int64()
}

func toInt32(i inter32) int64 {
	return i.Int64()
}

func toInt64(i inter64) int64 {
	return i.Int64()
}

func toIntStr(i interStr) int64 {
	return i.Int64()
}

func toIntPtr(i *interPtr) int64 {
	return i.Int64()
}

func toIntGeneric[T Inter](i T) int64 {
	return i.Int64()
}
