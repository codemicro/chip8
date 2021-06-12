package vm

// Stack is a unlimited size, 16-bit, FIFO stack
type Stack []uint16

func (s *Stack) Add(cont uint16) {
	*s = append(*s, cont)
}

func (s *Stack) Pop() uint16 {
	var x uint16
	x, *s = (*s)[len(*s)-1], (*s)[:len(*s)-1]
	return x
}