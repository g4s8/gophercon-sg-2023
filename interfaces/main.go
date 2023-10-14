package main

type fooer interface {
	foo()
}

type foo int

func (f foo) foo() {
	print("foo")
}

func main() {
	var f foo
	printFooer(f)
	printFoo(f)
}

func printFooer(f fooer) {
	f.foo()
}

func printFoo(f foo) {
	f.foo()
}

// ./main.go:19:17: parameter f leaks to {heap} with derefs=0:
// ./main.go:19:17:   flow: {heap} = f:
// ./main.go:19:17:     from f.foo() (call parameter) at ./main.go:20:7
// ./main.go:19:17: leaking param: f
// ./main.go:15:13: f escapes to heap:
// ./main.go:15:13:   flow: {heap} = &{storage for f}:
// ./main.go:15:13:     from f (spill) at ./main.go:15:13
// ./main.go:15:13:     from printFooer(f) (call parameter) at ./main.go:15:12
// ./main.go:15:13: f escapes to heap
