package main

import (
	"strconv"
)

func main() {
	var (
		two   UInt32 = 2
		three UInt32 = 3
		four  UInt32 = 4
		five  UInt32 = 5
	)
	var i UInt32
	i.Add(&i, &two)
	println(i.String())
	i.Mul(&i, &four)
	println(i.String())
	i.Sub(&i, &three)
	println(i.String())
	i.Div(&i, &five)
	println(i.String())
}

type UInt32 uint32

func (i *UInt32) Add(x, y *UInt32) *UInt32 {
	*i = *x + *y
	return i
}

func (i *UInt32) Sub(x, y *UInt32) *UInt32 {
	*i = *x - *y
	return i
}

func (i *UInt32) Mul(x, y *UInt32) *UInt32 {
	*i = *x * *y
	return i
}

func (i *UInt32) Div(x, y *UInt32) *UInt32 {
	*i = *x / *y
	return i
}

func (i *UInt32) String() string {
	return strconv.FormatUint(uint64(*i), 10)
}
