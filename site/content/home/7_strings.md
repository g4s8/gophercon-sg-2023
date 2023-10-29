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

Source code for printing "Hello, world!" from bytes.

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

{{%note%}}
And its assembly code.

`runtime.slicebytetostring` - copies data from slice to string.
From slice array to strings array.
To make it immutable.
{{%/note%}}

```x86asm{4}
LEAQ 0x30(SP), AX			
LEAQ 0x23(SP), BX			
MOVL $0xd, CX				
CALL runtime.slicebytetostring(SB)	
```

---

See "unsafe" package documentation (since go1.20):
 - `unsafe.SliceData(b)` - get pointer for underlying `b` array
 - `unsafe.String(ptr, l)` - returns `string` value with `b` pointer as backed
 bytes data pointer, and `l` as a string length (in bytes).

---

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

## Side effect

Mutable strings.

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
{{% /note %}}

---

## Concatenation

The fast method with 1 allocation.

```go{|4-7|4,8|9-12|8,13}
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

---

But it's the same as `strings.Builder` with `Grow()`.

---

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

---

The `String()` method of `strings.Builder`.

```go{}
func (b *Builder) String() string {
	return unsafe.String(unsafe.SliceData(b.buf), len(b.buf))
}
```

---

It is the same as `strings.Join()`.

---

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

{{%note%}}
Ignore separator.
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

{{%note%}}
Reusable buffers.
{{%/note%}}

