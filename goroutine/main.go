package main

type foo struct {
	x int
}

func main() {
	f := &foo{x: 1}
	go printFoo(f)
}

func printFoo(f *foo) {
	println(f.x)
}
