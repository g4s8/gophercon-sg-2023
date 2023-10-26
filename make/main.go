package main

func main() {
	var arr1 [1310720]int
	var arr2 [1310721]int
	fillSlice(arr1[:], 1, 2, 3, 4, 5)
	fillSlice(arr2[:], 1, 2, 3, 4, 5)

	s1 := make([]int, 100)
	fillSlice(s1, 1, 2, 3, 4, 5)
	print(s1[0])
	s2 := make([]int, 100000) // escapes to heap
	_ = s2[0]
	s3 := newSlice(1)
	_ = s3[0]

	_ = make([]int, 0, 100)
	_ = make([]int, 0, 10000000) // escapes to heap
	_ = newSliceCap(1)

	m1 := make(map[int]int, 100)
	m1[1] = 1
	print(m1[1])
	m2 := make(map[int]int, 1000000000000000000)
	m2[1] = 2
	print(m2[1])
	m3 := newMap(1)
	m3[1] = 3
	print(m3[1])

	_ = make(chan int, 100)
	_ = make(chan int, 10000000)
	_ = newChan(1)
}

func fillSlice(s []int, vals ...int) {
	for i, v := range vals {
		s[i] = v
	}
}

//go:noinline
func newSlice(size int) []int {
	return make([]int, size) // escapes to heap
}

//go:noinline
func newSliceCap(size int) []int {
	return make([]int, 0, size) // escapes to heap
}

//go:noinline
func newMap(size int) map[int]int {
	return make(map[int]int, size)
}

//go:noinline
func newChan(size int) chan int {
	return make(chan int, size)
}

//go:noinline
func useSlice(s []int) {
	_ = s
}
