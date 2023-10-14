package main

func main() {
	func() {
		var (
			val1 = 42 // val1 escapes to heap
			val2 = 42 // val2 escapes to heap
			val3 = 42 // val3 escapes to heap
		)
		s := make([]*int, 1)
		setSlicePtr(s, 0, &val1)

		m := make(map[int]*int)
		setMapPtr(m, 0, &val2)

		ch := make(chan *int, 1)
		sendChanPtr(ch, &val3)
	}()

	func() {
		var (
			val1 = 42
			val2 = 42
			val3 = 42
		)
		s := make([]int, 1)
		setSliceVal(s, 0, val1)

		m := make(map[int]int)
		setMapVal(m, 0, val2)

		ch := make(chan int, 1)
		sendChanVal(ch, val3)
	}()
}

func setSlicePtr(s []*int, i int, val *int) {
	s[i] = val
}

func sendChanPtr(ch chan *int, x *int) {
	ch <- x
}

func setMapPtr(m map[int]*int, i int, val *int) {
	m[i] = val
}

func setSliceVal(s []int, i int, val int) {
	s[i] = val
}

func sendChanVal(ch chan int, x int) {
	ch <- x
}

func setMapVal(m map[int]int, i int, val int) {
	m[i] = val
}

// # github.com/g4s8/memory-patterns/maps
// ./main.go:37:35: parameter val leaks to {heap} with derefs=0:
// ./main.go:37:35:   flow: {heap} = val:
// ./main.go:37:35:     from s[i] = val (assign) at ./main.go:38:7
// ./main.go:37:18: s does not escape
// ./main.go:37:35: leaking param: val
// ./main.go:45:39: parameter val leaks to {heap} with derefs=0:
// ./main.go:45:39:   flow: {heap} = val:
// ./main.go:45:39:     from m[i] = val (assign) at ./main.go:46:7
// ./main.go:45:16: m does not escape
// ./main.go:45:39: leaking param: val
// ./main.go:41:32: parameter x leaks to {heap} with derefs=0:
// ./main.go:41:32:   flow: {heap} = x:
// ./main.go:41:32:     from ch <- x (send) at ./main.go:42:5
// ./main.go:41:18: ch does not escape
// ./main.go:41:32: leaking param: x
// ./main.go:49:18: s does not escape
// ./main.go:57:16: m does not escape
// ./main.go:53:18: ch does not escape
// ./main.go:6:4: val1 escapes to heap:
// ./main.go:6:4:   flow: {heap} = &val1:
// ./main.go:6:4:     from &val1 (address-of) at ./main.go:11:21
// ./main.go:6:4:     from setSlicePtr(s, 0, &val1) (call parameter) at ./main.go:11:14
// ./main.go:7:4: val2 escapes to heap:
// ./main.go:7:4:   flow: {heap} = &val2:
// ./main.go:7:4:     from &val2 (address-of) at ./main.go:14:19
// ./main.go:7:4:     from setMapPtr(m, 0, &val2) (call parameter) at ./main.go:14:12
// ./main.go:8:4: val3 escapes to heap:
// ./main.go:8:4:   flow: {heap} = &val3:
// ./main.go:8:4:     from &val3 (address-of) at ./main.go:17:19
// ./main.go:8:4:     from sendChanPtr(ch, &val3) (call parameter) at ./main.go:17:14
// ./main.go:6:4: moved to heap: val1
// ./main.go:7:4: moved to heap: val2
// ./main.go:8:4: moved to heap: val3
// ./main.go:4:2: func literal does not escape
// ./main.go:10:12: make([]*int, 1) does not escape
// ./main.go:13:12: make(map[int]*int) does not escape
// ./main.go:20:2: func literal does not escape
// ./main.go:26:12: make([]int, 1) does not escape
// ./main.go:29:12: make(map[int]int) does not escape
