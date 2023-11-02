+++
weight = 10
+++

## Go low-latency patters

![qr-slides-repo](images/qr-slides-repo.png)

*Kirill Cherniavskiy for
GopherCon Singapore 2023*

---

# This talk

 - {{%fragment%}}See some code examples{{%/fragment%}}
 - {{%fragment%}}Analyze it with different tools{{%/fragment%}}
 - {{%fragment%}}Find problems and try to optimize{{%/fragment%}}

{{%note%}}
    In this  talk:
     - This talk is about code patterns and examples
     for low-latency code.
     - More examples and less theory.

Good to know

 - Go memory model.
 - Go internal types data structure.
 - A bit about GC, runtime, compiler.
{{%/note%}}

---

# Disclaimer

 - {{%fragment%}}Prefer readability where possible{{%/fragment%}}
 - {{%fragment%}}Prefer simple code, not complex{{%/fragment%}}
 - {{%fragment%}}Do not over-optimize without need{{%/fragment%}}
 - {{%fragment%}}Have a good reason and proof to optimize the code{{%/fragment%}}

{{%note%}}
Simple and readable code is better than
complex and unreadable.

Optimize only when sure about it and proved
by benchmarks or profiling.
{{%/note%}}

---

# Latency

> A time delay between the cause and the effect.

(c) Wikipedia

{{%note%}}
    - Qute: generally speaking
    - A time delay between event and reaction to the event.
    - Can't improve latency of hardware, OS,
    thread schedulers from code.
    - Can improve by removing random-time
    operations from code.
    - One of the main random fixable by code:
    GC and allocations, syncronization.
{{%/note%}}

---

## Common latency issues

 - {{%fragment%}}IO operations{{%/fragment%}}
 - {{%fragment%}}Syncronizations{{%/fragment%}}
 - {{%fragment%}}Heap allocations{{%/fragment%}}

{{%note%}}
 - Try using non blocking IO operations if possible.
 - Try to avoid syncronizations if possible,
 e.g. using atomics or single thread logic.
{{%/note%}}

---

## Real life allocation issue

Causes GC to eat **30% of CPU time**

```go{}
import "math/big"

type State struct {
    // skipping fields
    index *big.Int
}

func (s *state) UpdateIndex(val *big.Int) {
    s.index = new(big.Int).Set(val)
}
```

---

## The fix

```go{}
import "math/big"

type State struct {
    // skipping fields
    index big.Int
}

func (s *state) UpdateIndex(val *big.Int) {
    s.index.Set(val)
}
```

---

# Tools

Which tools I did use for this examples?


```sh{1-4|5}
$ go build -gcflags '-m' # simple escape analysis
$ go build -gcflags '-m=2' # more verbose analysis
$ go build -gcflags '-l' # disable inlining
$ go build -gcflags '-S' # print assembly listing
$ go tool objdump -s main.main -S example.com > main.go.s
```

{{%note%}}
 - different `gcflags`
 - objdump
 - benchmarks
 - profiler
{{%/note%}}

---

### lensm

[github.com/loov/lensm](https://github.com/loov/lensm)

![lens-screenshot](images/lensm.png)

---


# Content

 - Interfaces
 - Generics
 - Inlines
 - Pointers
