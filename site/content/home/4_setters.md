+++
weight = 32
+++

# Mutators

Assigning a pointer to a struct field causes heap allocation.

---

Example:
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

Both values are moved to heap:
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

Separate allocation and value assignment:
```go{2|5-7}
type foo struct {
	f *int
}

func (b *foo) setVal(v int) {
	*b.f = v
}
```

---

```go{2-3|5}
func main() {
	var target foo
	target.f = new(int) // call before critical path

	target.setVal(2) // call on performance critical path
}
```

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

## Dirty trick

```go{|1,4|7-9}
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
How to avoid it?
{{%/note%}}

---

```go{7-10|2}
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

**It could be dangerous** --- use only if the child object is not accessible outside of
the parent's stack frame.
