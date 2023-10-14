package main

type foo struct {
	x int
}

type bar struct {
	f *foo
}

func main() {
	var b1 bar
	b1.setDefault()
	print(&b1)

	var (
		b2 bar
		f  foo // f escapes to heap:
	)
	b2.setFoo(&f)
	print(&b2)

	var b3 bar
	b3.f = &foo{} // &foo{} does not escape
	b3.setFooX(1)
	print(&b3)
}

func (b *bar) setDefault() {
	b.f = &foo{} // &foo{} escapes to heap
}

func (b *bar) setFoo(f *foo) {
	b.f = f
}

func (b *bar) setFooX(x int) {
	if b.f == nil {
		panic("f is nil")
	}
	b.f.x = x
}

// # github.com/g4s8/memory-patterns/setters
// ./main.go:30:8: &foo{} escapes to heap:
// ./main.go:30:8:   flow: {heap} = &{storage for &foo{}}:
// ./main.go:30:8:     from &foo{} (spill) at ./main.go:30:8
// ./main.go:30:8:     from b.f = &foo{} (assign) at ./main.go:30:6
// ./main.go:29:7: b does not escape
// ./main.go:30:8: &foo{} escapes to heap
// ./main.go:33:22: parameter f leaks to {heap} with derefs=0:
// ./main.go:33:22:   flow: {heap} = f:
// ./main.go:33:22:     from b.f = f (assign) at ./main.go:34:6
// ./main.go:33:7: b does not escape
// ./main.go:33:22: leaking param: f
// ./main.go:39:9: "f is nil" escapes to heap:
// ./main.go:39:9:   flow: {heap} = &{storage for "f is nil"}:
// ./main.go:39:9:     from "f is nil" (spill) at ./main.go:39:9
// ./main.go:39:9:     from panic("f is nil") (call parameter) at ./main.go:39:8
// ./main.go:37:7: b does not escape
// ./main.go:39:9: "f is nil" escapes to heap
// ./main.go:18:3: f escapes to heap:
// ./main.go:18:3:   flow: {heap} = &f:
// ./main.go:18:3:     from &f (address-of) at ./main.go:20:12
// ./main.go:18:3:     from (*bar).setFoo(b2, &f) (call parameter) at ./main.go:20:11
// ./main.go:18:3: moved to heap: f
// ./main.go:24:9: &foo{} does not escape
