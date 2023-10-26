+++
weight = 40
+++

# Strings

What could be slow with strings in Go.

---

## Convert from bytes

Creating a new string from byte slice copies data to the new
string.

---

Source code for printing "Hello, world!":
```go{|7}
func main() {
	// bytes for "Hello, World!" string
	hello := []byte{
		72, 101, 108, 108, 111, 44, 32,
		87, 111, 114, 108, 100, 33,
	}
	helloStr := string(hello)
	println(helloStr)
}
```

---

And its assembly code:
```x86asm{4}
LEAQ 0x30(SP), AX			
LEAQ 0x23(SP), BX			
MOVL $0xd, CX				
CALL runtime.slicebytetostring(SB)	
```

---

`runtime.slicebytetostring` - copies data from slice to string.

{{%note%}}
From slice array to strings array.
To make it immutable.
{{%/note%}}

---

See "unsafe" package documentation (since go1.20):
 - `unsafe.SliceData(b)` - get pointer for underlying `b` array
 - `unsafe.String(ptr, l)` - returns `string` value with `b` pointer as backed
 bytes data pointer, and `l` as a string length (in bytes).

---

Using "unsafe" package to avoid copying data from byte slice:
```go{7-8}
func main() {
	// bytes for "Hello, World!" string
	hello := []byte{
		72, 101, 108, 108, 111, 44, 32,
		87, 111, 114, 108, 100, 33,
	}
	helloStr := unsafe.String(unsafe.SliceData(hello),
		len(hello))
	println(helloStr)
}
```

{{%note%}}
It's unsafe!
{{%/note%}}

---

Side effect --- mutable string.

```go{8-10}
func main() {
	hello := []byte{
		72, 101, 108, 108, 111, 44, 32,
		87, 111, 114, 108, 100, 33,
	}
	helloStr := unsafe.String(unsafe.SliceData(hello),
		len(hello))
	println(helloStr) // -> "Hello, World!"
	hello[0] = 'h'
	println(helloStr) // -> "hello, World!"
}
```

{{% note %}}
Remember this for the next topic: strings concatenation.

This leads to...
{{% /note %}}

---

## Concatenation
- `+=`
- `strings.Join()`
- `strings.Builder`
- `bytes.Buffer`
- `copy()`
- `unsafe`

---

Simple concatenation `+=`
```go{3-5}
func concatStringsAdd(ss []string) string {
	var s string
	for _, v := range ss {
		s += v
	}
	return s
}
```

---

`strings.Join()`
```go{4}
import "strings"

func concatStringsJoin(ss []string) string {
	return strings.Join(ss, "")
}
```

---

`strings.Builder`

```go{5-8}
import "strings"

func concatStringsBuilder(ss []string) string {
	var b strings.Builder
	for _, v := range ss {
		b.WriteString(v)
	}
	return b.String()
}
```

---

Improve with `Grow()`

```go{8,10-13|4-7,9}
import "strings"

func concatStringsBuilderGrow(ss []string) string {
    var size int
    for _, v := range ss {
        size += len(v)
    }
    var b strings.Builder
    b.Grow(size)
    for _, v := range ss {
        b.WriteString(v)
    }
    return b.String()
}
```

{{%note%}}
It was just appending string to buffer.
Now we calculate total length in bytes,
and Grow the buffer before writing.
{{%/note%}}

---

`bytes.Buffer`

```go{5-8}
import "bytes"

func concatStringsBuffer(ss []string) string {
	var b bytes.Buffer
	for _, v := range ss {
		b.WriteString(v)
	}
	return b.String()
}
```

---

With `Grow()`
```go{4-7,9}
import "bytes"

func concatStringsBufferGrow(ss []string) string {
    var size int
    for _, v := range ss {
        size += len(v)
    }
    var b bytes.Buffer
    b.Grow(size)
    for _, v := range ss {
        b.WriteString(v)
    }
    return b.String()
}
```

---

Just `copy()` bytes data

```go{2-6|7-10|11}
func concatStringsCopy(ss []string) string {
    var size int
    for _, v := range ss {
        size += len(v)
    }
    b := make([]byte, size)
    var i int
    for _, v := range ss {
        i += copy(b[i:], v)
    }
    return string(b)
}
```

---

Get rid of new allocation with `unsafe`:

```go{13}
import "unsafe"

func concatStringsCopyUnsafe(ss []string) string {
    var size int
    for _, v := range ss {
        size += len(v)
    }
    b := make([]byte, size)
    var i int
    for _, v := range ss {
        i += copy(b[i:], v)
    }
    return unsafe.String(unsafe.SliceData(b), len(b))
}
```

{{% note %}}
Remember casts and mutable string.
Here it's quite safe to use unsafe.
{{% /note %}}


---

### Benchmark results

Concat 100 strings of 100 bytes on AMD Ryzen 7 5700U.

{{%note%}}
Check all benchmarks on GitHub. here only one benchmark set.


// BuildeGrow, Join and CopyUnsafe - top (~ 1.9-2 mksec op on 100/100, just 1 alloc)
// BufferGrow and Copy - second (~3.5 mksec, 2 allocations)
// Just buffer and builder (~8.5 mksec, amount of allocations is logarithmic depends on amount of data)
// Concat - the worst (~90 mksec, 100 allocations)
{{%/note%}}

---

### Benchmark results (1st place)

1.9 - 2.0 μs and 1 allocation per operation
 - `strings.Builder` with `Grow`
 - `strings.Join`
 - unsafe `copy`

---

### Benchmark results (1st place)

Actually all has almost the same implementation:
 - `strings.Join` calls `strings.Builder` with `Grow`
 - `strings.Builder` uses `copy` and unsafe operations to build the string

---

Join implementation:
```go{4,5,6,8-10}
func Join(elems []string, sep string) string {
    // ...
    var b Builder
    b.Grow(n)
    b.WriteString(elems[0])
    for _, s := range elems[1:] {
        b.WriteString(sep)
        b.WriteString(s)
    }
    return b.String()
}
```

---

Builder implementation:

```go{2}
func (b *Builder) String() string {
	return unsafe.String(unsafe.SliceData(b.buf), len(b.buf))
}
```

---

### Other results

 - 3.5 μs and 2 allocations -- `bytes.Buffer` with `Grow` and `copy()`
 - 8.5 μs and log2(size) allocs -- `strings.Builder` and `bytes.Buffer`
 - 90 μs and 99 allocations -- `str += next`

{{%note%}}
Not a surprise - bytes.Buffer casts bytes slice to string, same as copy impl.
{{%/note%}}

---

To not reinvent the wheel just use `strings.Join` - it's the same
as a `string.Builder` with `Grow` and `copy` with "unsafe".

---

### Summary

 - "unsafe" helps to avoid copying data
But as a side effect, we get mutable string.
 - `strings.Join()` - fast and simple
 - `strings.Builder` with `Grow` - flexible and fast
 - `copy` with "unsafe" very flexible but the same as `Join()`


