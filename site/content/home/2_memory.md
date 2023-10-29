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

```go{1-3|5-8|10-14}
type Inter interface { // like fmt.Stringer
	Int64() int64
}

//go:noinline
func toInt(i Inter) int64 {
	return i.Int64()
}

type inter64 int64 // implementation

func (i inter64) Int64() int64 {
	return int64(i)
}
```

{{%note%}}
Inter interface is Like fmt.Stringer.

inter64 Int64() method has value receiver.
{{%/note%}}

---

## Does it move to heap?

```go{1-7}
func main() {
	i64_1 := inter64(1)
	_ = toInt(i64_1)

	i64_256 := inter64(256)
	_ = toInt(i64_256)
}

//go:noinline
func toInt(i Inter) int64 {
	return i.Int64()
}
```

---

Escape analysis reports that both values was moved to heap. Isn't it?

```txt{1,5,9,10}
i64_1 escapes to heap:
  flow: {heap} = &{storage for i64_1}:
    from i64_1 (spill) at ./main.go:50:12
    from toInt(i64_1) (call parameter) at ./main.go:50:11
i64_256 escapes to heap:
  flow: {heap} = &{storage for i64_256}:
    from i64_256 (spill) at ./main.go:55:12
    from toInt(i64_256) (call parameter) at ./main.go:55:11
i64_1 escapes to heap
i64_256 escapes to heap
```

---

The most popular fix.

```go{6-8|2-3}
func main() {
	i64_1 := inter64(1) // MOVL $0x1, AX
	_ = toInt64(i64_1)  // CALL main.toInt64(SB)
}

func toInt64(i inter64) int64 {
	return i.Int64()
}
```

{{%note%}}
One of possible solutions.

No allocations.
{{%/note%}}

---

## Does it really escape?

```go{1-7}
func main() {
	i64_1 := inter64(1)
	_ = toInt(i64_1)

	i64_256 := inter64(256)
	_ = toInt(i64_256)
}

//go:noinline
func toInt(i Inter) int64 {
	return i.Int64()
}
```

{{%note%}}
Question:
 - no escape
 - one variable is escaped
 - both escapes
{{%/note%}}

---

## Go deeper

```x86asm{1}
CALL runtime.convT64(SB)			
MOVQ AX, BX					
LEAQ go:itab.main.inter64,main.Inter(SB), AX	
CALL main.toInt(SB)	
```

{{%note%}}
The arguments for the function parameter is converted with `convT64`.
{{%/note%}}

---

The `runtime.convT64` function **may** moves value to heap
and returns pointer to the value.

---

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
{{%/note%}}

---

Where `staticuint64s` is an array of integers from `0x00` to `0xff`.

```go{}
// staticuint64s is used to avoid allocating in convTx
// for small integer values.
var staticuint64s = [...]uint64{
  0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
  // ...
  0xf8, 0xf9, 0xfa, 0xfb, 0xfc, 0xfd, 0xfe, 0xff,
}
```

---

```go{}
func main() {
	i64_1 := inter64(1) // doesn't escape
	_ = toInt(i64_1)    // get from statucuints64

	i64_256 := inter64(256) // escapes to heap
	_ = toInt(i64_256)      // allocate
}
```

---

## The takeaway

Escape analyzer output is just a hint for possible heap allocation:
verify it to be sure before starting to fix.

{{%note%}}
Escape analysis may be wrong, always use benchmarks for critical code.
{{%/note%}}

---

## Generic function

```go{6-8}
func main() {
	i64_1 := inter64(1)
	_ = toIntGeneric(i64_1)
}

func toIntGeneric[T Inter](i T) int64 {
	return i.Int64()
}
```

{{%note%}}
No allocation for `int64` here too.
{{%/note%}}

---

Compiler creates generic function for `go.shape.int64`,

{{%note%}}
Compiler generates generic functions for different GCShapes.
Pointer GCShape will be allocated on heap.
Generic function call can not be inlined.

The GC shape of a type means how that type appears
to the allocator / garbage collector.
It is determined by its size, its required alignment,
and which parts of the type contain a pointer.
{{%/note%}}

---

`inter64` argument is not moved to heap.
```x86asm{2}
LEAQ main..dict.toIntGeneric[main.inter64](SB), AX	
MOVL $0x1, BX						
CALL main.toIntGeneric[go.shape.int64](SB)
```

{{%note%}}
No allocation but additional overhead on types lookup in dict.
{{%/note%}}

---

Generic dynamic dispatch can't be inlined
```x86asm{1,3,6}
LEAQ main..dict.toIntGeneric[main.inter64](SB), AX	
; ---
MOVQ 0(AX), CX		
MOVQ AX, DX		
MOVL BX, AX		
CALL CX			
ADDQ $0x8, SP		
```

{{%note%}}
optimization that bypasses v-table lookup for inline functions.
{{%/note%}}

---

## Inlines

If function with interface parameter is inlined, compiler may not move to heap its
arguments.

{{%note%}}
And it usually has better performance than generic.
Because generic function can't be inlined and has dynamic calls.
All previous examples I've compiled with disable inlines gcflags option.
{{%/note%}}

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
Another example
{{%/note%}}

---

Let's see the difference between interface parameter and generic
parameter.

---

```go{1-6|8-13}
func sum(calc Calc, vals ...int) (sum int) {
	for _, val := range vals {
		sum = calc.Add(val)
	}
	return
}

func sumGeneric[T Calc](calc T, vals ...int) (sum int) {
	for _, val := range vals {
		sum = calc.Add(val)
	}
	return
}
```
---

```go{}
// func sum(calc Calc, vals ...int) (sum int)

func main() {
	var c1 calcInt
	_ = sum(&c1, 1, 2, 3)
}
```

{{%note%}}
Calling function with interface parameter.
{{%/note%}}

---

Inlined implementation call.

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

Generic function.
{{%note%}}
Calling function with generic parameter type:
{{%/note%}}

```go{}
// func sumGeneric[T Calc](calc T, vals ...int) (sum int)

func main() {
	var c2 calcInt
	_ = sumGeneric(&c2, 1, 2, 3)
}
```

---

No `Add()` inline:
```x86asm{1-4|7-10}
; var c1g calcInt
LEAQ 0x4b86(IP), AX		
CALL runtime.newobject(SB)	
MOVQ AX, 0x90(SP)		
; _ = sumGeneric(&c1g, 1, 2, 3)
; ... range loop jumps
sum = calc.Add(val)
LEAQ main..dict.sumGeneric[*main.calcInt](SB), DX	
LEAQ 0xfffffef6(IP), SI					
CALL SI
```

---

## Pros and cons

**Actual type parameters**
 - {{%fragment%}}Can avoid allocation ➕{{%/fragment%}}
 - {{%fragment%}}No dynamic calls ➕{{%/fragment%}}
 - {{%fragment%}}Could be inlined ➕{{%/fragment%}}
 - {{%fragment%}}Write more code ➖{{%/fragment%}}

---

## Pros and cons

**Interface type parameters**
 - {{%fragment%}}Could be inlined ➕{{%/fragment%}}
 - {{%fragment%}}More readable code ➕{{%/fragment%}}
 - {{%fragment%}}Often moved to heap if not inlined ➖{{%/fragment%}}
 - {{%fragment%}}Dynamic dispatch if not inlined ➖{{%/fragment%}}

---

## Pros and cons

**Generic type parameters**
 - {{%fragment%}}More readable code ➕{{%/fragment%}}
 - {{%fragment%}}No allocation for some GCShapes ➕{{%/fragment%}}
 - {{%fragment%}}Dynamic dispatch ➖{{%/fragment%}}
 - {{%fragment%}}Call not be inlined ➖{{%/fragment%}}

---

## Summary

 - Check if interface could be inlined
 - Try using generic method
 - Use `func` with actual types
 - Or redesign the function

{{%note%}}
 - Quite often small functions with interface parameters could be inlined.

 - Generic function may not allocate memory for passing types for some GCShapes,
 but remember about dynamic dispatch.

 - If there are no so many different types for function parameters, try to
 implement different functions for each type. It may improve the performance
 comparing with interface parameter function but also it can make it worser.

 - Try to think about changing data type to generic instead of using generic
 or interface function parameters types.

{{%/note%}}

