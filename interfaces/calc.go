package main

type Calc interface {
	Add(int) int
}

type calcInt int

func (c *calcInt) Add(n int) (sum int) {
	sum = int(*c) + int(n)
	*c = calcInt(sum)
	return
}

type calcStruct struct {
	sum int
}

func (c *calcStruct) Add(n int) (sum int) {
	sum = c.sum + n
	c.sum = sum
	return
}

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
