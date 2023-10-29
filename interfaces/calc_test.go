package main

import "testing"

func TestCalcInt(t *testing.T) {
	var c calcInt
	if sum(&c, 1, 2, 3) != 6 {
		t.Fail()
	}
	if sumGeneric(&c, 1, 2, 3, 4) != 16 {
		t.Fail()
	}
	if c != 16 {
		t.Fail()
	}
}
