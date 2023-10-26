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

### Examples

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

Which value is moved to the heap?
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
Ask which argument is moved to heap:
 - none
 - one of them
 - both
{{%/note%}}

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

## Go deeper

The arguments for the function parameter is converted with `convT64`:

```x86asm{1}
CALL runtime.convT64(SB)			
MOVQ AX, BX					
LEAQ go:itab.main.inter64,main.Inter(SB), AX	
CALL main.toInt(SB)	
```

{{%note%}}
In both cases convT64
{{%/note%}}

---

The `runtime.convT64` function **may** move value to heap
and returns pointer to the value:
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

Where `staticuint64s` is an array of integers from `0x00` to `0xff`:
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

### The bottom line

Escape analyzer report is just a hint for possible heap allocation:
check it to be sure.

{{%note%}}
Escape analysis may be wrong, always use benchmarks for critical code.

BTW, calling a function with interface parameters leads to types lookup in itable.
{{%/note%}}

---

### Using exact parameters

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

### Generic function

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
where the argument is passed by value.

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

It has dynamic call → can't be inlined:
```x86asm{1,3,6}
LEAQ main..dict.toIntGeneric[main.inter32](SB), AX	
; ---
MOVQ 0(AX), CX		
MOVQ AX, DX		
MOVL BX, AX		
CALL CX			
ADDQ $0x8, SP		
```

---

### Pointer type

```go{}
type interPtr int64

func (i *interPtr) Int64() int64 {
	return int64(*i)
}
```

---

Pointer type arguments are moved to heap for interface and generic function parameters.

---

```go{2,3,7-10|2,4,11-14}
func main() {
	iPtr := interPtr(1)
	_ = toInt(&iPtr)
	_ = toIntGeneric(&iPtr)
}

//go:noinline
func toInt(i Inter) int64 {
	return i.Int64()
}

func toIntGeneric[T Inter](i T) int64 {
	return i.Int64()
}
```

{{%note%}}
Moved to heap and dynamic call.
{{%/note%}}

---

No allocation for exact types:

```go{}
func main() {
	iPtr := interPtr(1)
	_ = toIntPtr(&iPtr)
}

func toIntPtr(i *interPtr) int64 {
	return i.Int64()
}
```

{{%note%}}
More examples in repository
{{%/note%}}

---

More examples in repository:
 - `int32`
 - `string`
 - `type foo struct`

---

If function with interface parameter is inlined, compiler may not move to heap its
arguments.

{{%note%}}
And it usually has better performance than generic.
Because generic function can't be inlined and has dynamic calls.
All previous examples I've compiled with disable inlines gcflags option.
{{%/note%}}

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
 - {{%fragment%}}Dynamic call if not inlined ➖{{%/fragment%}}

---

## Pros and cons

**Generic type parameters**
 - {{%fragment%}}More readable code ➕{{%/fragment%}}
 - {{%fragment%}}No allocation for some GCShapes ➕{{%/fragment%}}
 - {{%fragment%}}Dynamic call ➖{{%/fragment%}}
 - {{%fragment%}}Can not be inlined ➖{{%/fragment%}}

---

## Summary

 - Check if interface could be inlined
 - Try generic method for primitive types
 - Use `func` with actual types
 - Or redesign the function

{{%note%}}
 - Quite often small functions with interface parameters could be inlined.

 - If there are no so many different types for function parameters, try to
 implement different functions for each type.

 - Generic method may have better performance for primitive types.
 Try to think about changing data type to generic instead of using generic
 function parameters.

{{%/note%}}
