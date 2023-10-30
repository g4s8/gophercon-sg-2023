+++
weight = 10
+++

## Go low-latency patters

![qr-slides-repo](images/qr-slides-repo.png)

*Kirill Cherniavskiy for
GopherCon Singapore 2023*

---

# This talk

 - See some code examples
 - Analyze it with different tools
 - Find problems and try to optimize

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

# Latency

> A time delay between the cause and the effect.

(c) Wikipedia

{{%note%}}
    - Qute: generally speaking
    - A time delay between event and reaction to the event.
    - Can't improve latency of hardware, OS, thread schedulers
    from code.
    - Can improve by removing random-time operations from code.
    - One of the main random fixable by code: GC and allocations,
    syncronization.
{{%/note%}}

---

# Disclaimer

 - {{%fragment%}}Prefer readability where possible.{{%/fragment%}}
 - {{%fragment%}}Prefer simple code, not complex.{{%/fragment%}}
 - {{%fragment%}}Do not over optimize without need.{{%/fragment%}}
 - {{%fragment%}}Have a good reason and proof to optimize the code.{{%/fragment%}}

{{%note%}}
Simple and readable code is better than
complex and unreadable.

Optimize only when sure about it and proved
by benchmarks or profiling.
{{%/note%}}

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
 - Mutators
