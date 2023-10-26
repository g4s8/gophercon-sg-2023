package main

import (
	"strconv"
)

func main() {
	one := new(SmallInt).Set(1)
	two := new(SmallInt).Set(2)
	three := new(SmallInt).Set(3)

	var sum SmallInt
	sum.Add(&sum, three).Add(&sum, two).Sub(&sum, one)
	println("3 + 2 - 1 = ", sum.String())
}

type SmallInt [1]int32

func (i *SmallInt) Set(x int32) *SmallInt {
	i[0] = x
	return i
}

func (i *SmallInt) Add(x, y *SmallInt) *SmallInt {
	i[0] = x[0] + y[0]
	return i
}

func (i *SmallInt) Sub(x, y *SmallInt) *SmallInt {
	i[0] = x[0] - y[0]
	return i
}

func (i *SmallInt) String() string {
	return strconv.FormatInt(int64(i[0]), 10)
}
