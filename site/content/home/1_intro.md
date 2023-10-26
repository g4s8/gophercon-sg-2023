+++
weight = 10
+++

## Go low-latency patters

![qr-slides-repo](images/qr-slides-repo.png)

*Kirill Cherniavskiy for
GopherCon Singapore 2023*

---

## This talk

 - See some code examples
 - Analyze it with different tools
 - Check assembly code
 - Find problems and try to optimize

{{%note%}}
    In this  talk:
     - This talk is about code patterns and examples
     for low-latency code.
     - More examples and less theory.
     - It's not about algorithms, it's about Go internals.

    - It could be useful when your code should react to external events
    in a fast and predictable time.
    - Of course, there will be some unpredictable factors,
    like hardware, OS scheduler, etc.
{{%/note%}}

---

# Disclaimer

 - Prefer readability where possible.
 - Prefer simple code, not complex.
 - Do not over optimize without need.

{{%note%}}
Simple and readable code is better than
complex and unreadable.

Optimize only when sure about it and proved
by benchmarks or profiling.
{{%/note%}}

---

### Good to know

 - Go memory model
 - Go internal types data structure
 - Garbage collector
 - A bit of runtime and compiler

{{%note%}}
Allocating memory in a heap may affect latency because of GC.

Object is moved to heap if compiler can't prove that
it's not accessible after function return.
{{%/note%}}

---

### Real life story

This one-line function executed a few millions times/sec
caused GC to eat about 30% of CPU:
```go{5-7}
type container struct {
	index *big.Int
}

func (c *container) updateIndex(val *big.Int) {
	c.index = new(big.Int).Set(val)
}
```

---

Fix:
```go{2,6}
type container struct {
	index big.Int
}

func (c *container) updateIndex(val *big.Int) {
	c.index.Set(val)
}
```

---

# Tools

Which tools I did use for this examples?

---

### go build

With gcflags:

```sh{1|2|3|4}
$ go build -gcflags '-m' # simple escape analysis
$ go build -gcflags '-m=2' # more verbose analysis
$ go build -gcflags '-l' # disable inlining
$ go build -gcflags '-S' # print assembly listing
```

---

### assembly
```bash
$ go tool objdump -s main.main -S example.com > main.go.s
```

---

### lensm

[github.com/loov/lensm](https://github.com/loov/lensm)

![lens-screenshot](images/lensm.png)

---

### Benchmarks

```go{}
package main

import "testing"

func BenchmarkCaseOne(b *testing.B) {
    // ...
}
```

---

## Content

 - Interfaces and generics
 - Mutators
 - Strings
