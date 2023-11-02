+++
weight = 32
+++

# Pointers

 - {{%fragment%}}Assigning a pointer to a struct field{{%/fragment%}}
 - {{%fragment%}}Returning a pointer{{%/fragment%}}

{{%note%}}
Two primary problems with pointers.
{{%/note%}}

---

## Assign to the field
```go{|5-8|10-12}
type foo struct {
	f *int
}

func (b *foo) setDefault() {
	var one int = 1 // moved to heap: one
	b.f = &one
}

func (b *foo) setF(f *int) {
	b.f = f
}
```

---

Both values are moved to heap.

```go{2-3|5-7}
func main() {
	var b1 foo
	b1.setDefault()

	var b2 foo
	var f int = 2 // moved to heap: f
	b2.setF(&f)
}
```

---

### Separate alloc and assignment

```go{2,3,9|5,12-13}
func main() {
	var target foo
	target.f = new(int) // call before critical path

	target.setVal(2) // call on performance critical path
}

type foo struct {
	f *int
}

func (b *foo) setVal(v int) {
	*b.f = v
}
```

---

If you prefer pointers parameters.

```go{}
type foo struct {
	f *int
}

func (b *foo) setValPtr(v *int) {
	*b.f = *v
}
```

---

## Returns

```go{}
type foo struct {
    x int
}

func newFoo(x int) *foo {
    return &foo{x} // escapes to heap
}
```

---

### Set the value

```go{}
type foo struct {
    x int
}

func (f *foo) set(x int) *foo {
    f.x = x
    return f
}

func main() {
    f := new(foo).set(42) // no allocation
}
```

---

## Returning a pointer

 - {{%fragment%}}Parameter pointer{{%/fragment%}}
 - {{%fragment%}}Method receiver pointer{{%/fragment%}}

---

Types in `math/big` are a good example of a design that avoid redundant allocations.

---

No allocations:
```go{|4-6|7-9}
import "math/big"

func main() {
	one := new(big.Int).SetInt64(1)
	two := new(big.Int).SetInt64(2)
	three := new(big.Int).SetInt64(3)
	var sum big.Int
	sum.Add(&sum, one).Add(&sum, two).Add(&sum, three)
	println(sum.String())
}
```

{{%note%}}
How to design similar types?
 - Do not store pointers in type struct
 - Return only the pointers which were passed as params
{{%/note%}}

---

### How to create similar types?

```go{|1|3-6|8-11}
type SmallInt [1]int32

func (i *SmallInt) Set(x int32) *SmallInt {
	i[0] = x
	return i
}

func (i *SmallInt) Add(x, y *SmallInt) *SmallInt {
	i[0] = x[0] + y[0]
	return i
}
```

---

## How to bypass allocation

```go{}
type Child int

type Parent struct {
	C *Child
}

func (p *Parent) SetChild(c *Child) {
	p.C = c
}
```

---

```x86asm{1,3,4,5,9}
; c := Child(1)
LEAQ 0x4ef9(IP), AX		
CALL runtime.newobject(SB)	
MOVQ $0x1, 0(AX)		
; p.SetChild(&c)
MOVQ AX, BX				
LEAQ 0x20(SP), AX			
NOPL 0(AX)(AX*1)			
CALL main.(*Parent).SetChild(SB)	
```

{{%note%}}
As expected it's moved to heap.
{{%/note%}}

---

## Dirty hack

```go{5-10|2}
func (p *Parent) SetChildUnsafe(c *Child) {
	p.C = (*Child)(noescape(unsafe.Pointer(c)))
}

//go:nosplit
//go:nocheckptr
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}
```

{{%note%}}
noescape hides a pointer from escape analysis.
it breaks the dependency between parameter and returned value.
noescape can be inlined and compiles down to zero instructions.

Got it from Go runtime sources.

{{%/note%}}

---

```x86asm{1,2,3,6}
; c := Child(2)
MOVQ $0x2, 0x10(SP)	
; p.SetChildUnsafe(&c)
LEAQ 0x18(SP), AX			
LEAQ 0x10(SP), BX			
CALL main.(*Parent).SetChildUnsafe(SB)	
```

---

## Warning

**It could be dangerous** --- use only if the child object is not accessible outside of
the parent's stack frame.

---

### Cleanup after this

```go{}
func dangerousOperation(p *Parent) {
    defer p.SetChild(nil)
    var c Child
    p.SetChildUnsafe(&c)
    // work with parent and child
}
```
