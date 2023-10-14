package main

type foo struct {
	x int
}

func main() {
	f1 := newFoo(1)
	f2 := newFooStruct(2)
	var f foo
	f3 := makeFoo(&f, 3)
	print(f1, &f2, f3)
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

// ./main.go:16:9: &foo{...} escapes to heap:
// ./main.go:16:9:   flow: ~r0 = &{storage for &foo{...}}:
// ./main.go:16:9:     from &foo{...} (spill) at ./main.go:16:9
// ./main.go:16:9:     from return &foo{...} (return) at ./main.go:16:2
// ./main.go:16:9: &foo{...} escapes to heap
// ./main.go:23:14: parameter f leaks to ~r0 with derefs=0:
// ./main.go:23:14:   flow: ~r0 = f:
// ./main.go:23:14:     from return f (return) at ./main.go:25:2
// ./main.go:23:14: leaking param: f to result ~r0 level=0
