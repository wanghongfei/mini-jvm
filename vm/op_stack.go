package vm

// 操作数栈
type OpStack struct {
	elems []uint32

	// 永远指向栈顶元素
	topIndex int
}

func NewOpStack(maxDepth int) *OpStack {
	return &OpStack{
		elems:        make([]uint32, maxDepth),
		topIndex:    -1,
	}
}

// 压栈
func (s *OpStack) Push(data uint32) bool {
	if s.topIndex == len(s.elems) - 1 {
		// 栈满了
		return false
	}

	s.topIndex++
	s.elems[s.topIndex] = data

	return true
}

// 出栈
func (s *OpStack) Pop() (uint32, bool) {
	if -1 == s.topIndex {
		// 栈空
		return 0, false
	}

	data := s.elems[s.topIndex]
	s.topIndex--

	return data, true
}

