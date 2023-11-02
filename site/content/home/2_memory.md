+++
weight = 30
+++

# Interfaces

Be careful with interface function parameters --- arguments for these parameters
are often moved to heap before passing to the callee function.

{{%note%}}
Actually it's not always obvious when the value will be moved to heap and when it's not.
{{%/note%}}

---

{{< slide class="interface-structure" >}}

## Interface structure


{{< mermaid >}}
graph TD
  param-->*value
  param-->itable
  itable-->type
  itable-->funcs
{{< /mermaid >}}

{{%note%}}
Event if we pass value argument to function,
compiler will copy the data and get the pointer
to construct interface data type.

Basically, when the compiler can't prove it will
not be accessed after function return.
{{%/note%}}

---

Minimal reproducible heap escape.

```go{}
package main

import "fmt"

func main() {
	var x int = 256
	fmt.Println(x)
}
```

{{%note%}}
Check with escape aanalysis.
Check with benchmarks.
Try different x values.
{{%/note%}}

---

### Sometimes it's not obvoius when the argument escapes to heap

---

## Why 256?

```go{}
package main

import "fmt"

func main() {
	var x1, x2 int = 1, 256 // 0x01, 0xFF+1
	fmt.Println(x1)
	fmt.Println(x2)
}
```

---

Escape analysis output.

```txt{2,4}
./main.go:7:13: ... argument does not escape
./main.go:7:14: x1 escapes to heap
./main.go:8:13: ... argument does not escape
./main.go:8:14: x2 escapes to heap
```

---

Output of `objdump`.

```x86asm{1,2,10}
MOVL $0x1, AX			
CALL runtime.convT64(SB)	
LEAQ 0x6e9b(IP), CX		
MOVQ CX, 0x18(SP)		
MOVQ AX, 0x20(SP)		
LEAQ 0x18(SP), AX		
MOVL $0x1, BX			
MOVQ BX, CX			
NOPL 0(AX)			
CALL fmt.Println(SB)	
```

---

## Go deeper

```go{|3|5,6}
func convT64(val uint64) (x unsafe.Pointer) {
	if val < uint64(len(staticuint64s)) {
		x = unsafe.Pointer(&staticuint64s[val])
	} else {
		x = mallocgc(8, uint64Type, false)
		*(*uint64)(x) = val
	}
	return
}
```

{{%note%}}
Two branch:
 1. when value is less then cache size
 2. otherwise allocate memory

Escape analysis may be wrong, always use benchmarks for critical code.

Escape analyzer output is just a hint for possible heap allocation:
verify it to be sure before starting to fix.
{{%/note%}}

---

Only `x2` is allocated on heap.

```go{}
package main

import "fmt"

func main() {
	var x1, x2 int = 1, 256 // 0x01, 0xFF+1
	fmt.Println(x1)
	fmt.Println(x2)
}
```

---

### Putting a word-sized-or-less non-pointer type in an interface value doesn't allocate

---

```go{}
package main

import "fmt"

func main() {
	fmt.Println(1)
	var one int = 1
	fmt.Println(one)
}
```


---

### Println(1)

```x86asm{4,5,6}
MOVUPS X15, 0x18(SP)	
LEAQ 0x6ea5(IP), DX	
MOVQ DX, 0x18(SP)	
LEAQ 0x36f71(IP), SI	
MOVQ SI, 0x20(SP)	
LEAQ 0x18(SP), AX	
MOVL $0x1, BX		
MOVQ BX, CX		
CALL fmt.Println(SB)	
```

---

### Println(one)

```x86asm{2,3,6,7}
MOVUPS X15, 0x18(SP)		
MOVL $0x1, AX			
CALL runtime.convT64(SB)	
LEAQ 0x6e6b(IP), DX		
MOVQ DX, 0x18(SP)		
MOVQ AX, 0x20(SP)		
LEAQ 0x18(SP), AX		
MOVL $0x1, BX			
MOVQ BX, CX			
CALL fmt.Println(SB)	
```

---

## How to avoid allocations?

(if we have it and if we need it)

{{%note%}}
Allocations for interface parameters
{{%/note%}}

---

### Interface parameter example

```go{1-3|5-9|11-14}
type Inter interface { // like fmt.Stringer
	Int64() int64
}

type inter64 int64 // implementation

func (i inter64) Int64() int64 {
	return int64(i)
}

//go:noinline
func toInt(i Inter) int64 {
	return i.Int64()
}
```

{{%note%}}
Inter interface is Like fmt.Stringer.

inter64 Int64() method has value receiver.
{{%/note%}}

---



### Exact type parameter

```go{11-13}
type Inter interface {
	Int64() int64
}

type inter64 int64

func (i inter64) Int64() int64 {
	return int64(i)
}

//go:noinline
func toInt64(i inter64) int64 {
	return i.Int64()
}
```

{{%note%}}
One of possible solutions.

No allocations.

Bad solution from the design perspective.
{{%/note%}}

---

### Argument in register

```go{}
func main() {
	i64_1 := inter64(1) // MOVL $0x1, AX
	_ = toInt64(i64_1)  // CALL main.toInt64(SB)
}
```
---

## Generic function

```go{11-13}
type Inter interface {
	Int64() int64
}

type inter64 int64

func (i inter64) Int64() int64 {
	return int64(i)
}

//go:noinline
func toIntGeneric[T Inter](i T) int64 {
	return i.Int64()
}
```

---

### No escape for int64 shape

```x86asm{2}
LEAQ main..dict.toIntGeneric[main.inter64](SB), AX	
MOVL $0x1, BX						
CALL main.toIntGeneric[go.shape.int64](SB)
```

{{%note%}}
No allocation.
Compiler creates generic function for `go.shape.int64`:
{{%/note%}}

---


{{%note%}}
Compiler generates generic functions for different GCShapes.
Pointer GCShape will be allocated on heap.
Generic function call can not be inlined.


The GC shape of a type means how that type appears
to the allocator / garbage collector.
It is determined by its size, its required alignment,
and which parts of the type contain a pointer.
{{%/note%}}


## GC shape

> The GC shape of a type means how that type appears to the allocator / garbage collector.
> It is determined by its size, its required alignment,
> and which parts of the type contain a pointer.

---


## Inline optimization

When function call with interface parameter is inlined, it bypasses virtual table lookup.
Compiler may not move to heap its arguments.

{{%note%}}
In interface it's optimization that bypasses v-table lookup for inline functions.
All previous examples I've compiled with disable inlines gcflags option.


It usually has better performance than generic.
Because generic function can't be inlined and has dynamic calls.
{{%/note%}}

---

{{%note%}}
https://github.com/golang/go/wiki/CompilerOptimizations
{{%/note%}}

## Inline rules 

 - the number of AST nodes must less than 80;
 - doesn't contain closures, defer, recover, select;
 - isn't prefixed by `go:noinline`;
 - isn't prefixed by `go:uintptrescapes`;
 - function has body;

---

```go{1-3|5-11}
type Calc interface {
	Add(int) int
}

type calcInt int

func (c *calcInt) Add(n int) (sum int) {
	sum = int(*c) + int(n)
	*c = calcInt(sum)
	return
}
```

{{%note%}}
Another example: interface with method to call it.
{{%/note%}}

---

Let's see the difference between **interface parameter** vs **generic
parameter**.

---

### Interface parameter

```go{6-11|1-4}
func main() {
	var c1 calcInt
	_ = sum(&c1, 1, 2, 3)
}

func sum(calc Calc, vals ...int) (sum int) {
	for _, val := range vals {
		sum = calc.Add(val)
	}
	return
}
```

---

### Implementation call was inlined

```x86asm{1-2|3-8|9|10-12}
; var c1 calcInt
MOVQ $0x0, 0x18(SP)	
; _ = sum(&c1, 1, 2, 3)
MOVUPS X15, 0x70(SP)	
MOVUPS X15, 0x78(SP)	
MOVQ $0x1, 0x70(SP)	
MOVQ $0x2, 0x78(SP)	
MOVQ $0x3, 0x80(SP)	
; ... range loop jumps
; sum = calc.Add(val)
LEAQ 0x18(SP), AX		
CALL main.(*calcInt).Add(SB)
; ... the rest
```

---

### Generic function

```go{6-11|1-4}
func main() {
	var c2 calcInt
	_ = sumGeneric(&c2, 1, 2, 3)
}

func sumGeneric[T Calc](calc T, vals ...int) (sum int) {
	for _, val := range vals {
		sum = calc.Add(val)
	}
	return
}
```

---

### No inline for generic func

```x86asm{1-4|7-10}
; var c1g calcInt
LEAQ 0x4b86(IP), AX		
CALL runtime.newobject(SB)	
MOVQ AX, 0x90(SP)		
; _ = sumGeneric(&c1g, 1, 2, 3)
; ... range loop jumps
; sum = calc.Add(val)
LEAQ main..dict.sumGeneric[*main.calcInt](SB), DX	
LEAQ 0xfffffef6(IP), SI					
CALL SI
```

---

## Pros and cons


| Type         | Allocs        | Dyn. dispatch    | Inline     |
|--------------|---------------|------------------|------------|
| Exact        | No            | No               | Yes        |
| Interface    | Inline        | Inline           | Yes        |
| Generics     | Shapes        | Yes              | No         |

---

## Summary

 - {{%fragment%}}Check if interface could be inlined{{%/fragment%}}
 - {{%fragment%}}Maybe another design?{{%/fragment%}}
 - {{%fragment%}}Try using generic method{{%/fragment%}}
 - {{%fragment%}}Use `func` with actual types{{%/fragment%}}

{{%note%}}
 - Quite often small functions with interface parameters could be inlined.

 - Try to think about changing data type to generic instead of using generic
 or interface function parameters types.

 - Generic function may not allocate memory for passing types for some GCShapes,
 but remember about dynamic dispatch.

 - If there are no so many different types for function parameters, try to
 implement different functions for each type. It may improve the performance
 comparing with interface parameter function but also it can make it worser.
{{%/note%}}

