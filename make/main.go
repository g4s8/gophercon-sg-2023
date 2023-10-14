package main

func main() {
	_ = make([]int, 100)
	_ = make([]int, 100000) // escapes to heap
	newSlice(1)

	_ = make(map[int]int, 100)
	_ = make(map[int]int, 10000000)
	newMap(1)
}

func newSlice(size int) {
	_ = make([]int, size) // escapes to heap
}

func newMap(size int) {
	_ = make(map[int]int, size) // escapes to heap
}

// # github.com/g4s8/memory-patterns/make
// ./main.go:14:10: make([]int, size) escapes to heap:
// ./main.go:14:10:   flow: {heap} = &{storage for make([]int, size)}:
// ./main.go:14:10:     from make([]int, size) (non-constant size) at ./main.go:14:10
// ./main.go:14:10: make([]int, size) escapes to heap
// ./main.go:18:10: make(map[int]int, size) does not escape
// ./main.go:5:10: make([]int, 100000) escapes to heap:
// ./main.go:5:10:   flow: {heap} = &{storage for make([]int, 100000)}:
// ./main.go:5:10:     from make([]int, 100000) (too large for stack) at ./main.go:5:10
// ./main.go:4:10: make([]int, 100) does not escape
// ./main.go:5:10: make([]int, 100000) escapes to heap
// ./main.go:8:10: make(map[int]int, 100) does not escape
// ./main.go:9:10: make(map[int]int, 10000000) does not escape
