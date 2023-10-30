package main

import "unsafe"

type foo struct {
	f *int
}

func main() {
	var b1 foo
	b1.setDefault()

	var b2 foo
	var f int = 2
	b2.setF(&f)

	var b3 foo
	b3.f = new(int)
	b3.setVal(3)

	var b4 foo
	var f4 int = 4
	b4.SetUnsafe(&f4)

	var b5 foo
	var f5 int = 5
	b5.setValPtr(&f5)
}

func (b *foo) setDefault() {
	var z int = 1
	b.f = &z
}

func (b *foo) setF(f *int) {
	b.f = f
}

func (b *foo) setVal(v int) {
	*b.f = v
}

func (b *foo) setValPtr(v *int) {
	*b.f = *v
}

func (b *foo) SetUnsafe(c *int) {
	b.f = (*int)(noescape(unsafe.Pointer(c)))
}

//go:nosplit
//go:nocheckptr
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}
