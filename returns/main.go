package main

import "strconv"

type foo struct {
	x int
}

func main() {
	f1 := newFoo(1)
	f2 := newFooStruct(2)
	f3 := makeFoo(new(foo), 3)
	f4 := new(foo).set(4)
	println(f1.String(), f2.String(), f3.String(), f4.String())
}

func newFoo(x int) *foo {
	return &foo{x}
}

func newFooStruct(x int) foo {
	return foo{x}
}

func makeFoo(f *foo, x int) *foo {
	f.x = x
	return f
}

func (f *foo) set(x int) *foo {
	f.x = x
	return f
}

func (f *foo) String() string {
	return strconv.Itoa(f.x)
}
