package vm

import (
	"fmt"
	"testing"
)

func TestOpStack_Push(t *testing.T) {
	s := NewOpStack(5)

	s.Push(1)
	s.Push(2)
	s.Push(3)
	s.Push(4)
	s.Push(5)

	fmt.Println(s.Push(100))

	fmt.Println(s.Pop())
	fmt.Println(s.Pop())
	fmt.Println(s.Pop())
	fmt.Println(s.Pop())
	fmt.Println(s.Pop())

	fmt.Println(s.Pop())
}
