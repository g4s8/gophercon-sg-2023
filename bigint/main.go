package main

import "math/big"

func main() {
	one := new(big.Int).SetInt64(1)
	two := new(big.Int).SetInt64(2)
	three := new(big.Int).SetInt64(3)
	var sum big.Int
	sum.Add(&sum, one).Add(&sum, two).Add(&sum, three)
	println(sum.String())
}
