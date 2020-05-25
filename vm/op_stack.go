package vm

import "github.com/wanghongfei/mini-jvm/vm/class"

// 操作数栈
type OpStack struct {
	elems []interface{}

	// 永远指向栈顶元素
	topIndex int
}


func NewOpStack(maxDepth int) *OpStack {
	return &OpStack{
		elems:        make([]interface{}, maxDepth),
		topIndex:    -1,
	}
}

// 压栈
func (s *OpStack) Push(data interface{}) bool {
	if s.topIndex == len(s.elems) - 1 {
		// 栈满了
		return false
	}

	s.topIndex++
	s.elems[s.topIndex] = data

	return true
}

// 出栈
func (s *OpStack) Pop() (interface{}, bool) {
	if -1 == s.topIndex {
		// 栈空
		return nil, false
	}

	data := s.elems[s.topIndex]
	s.topIndex--

	return data, true
}

func (s *OpStack) GetTop() (interface{}, bool) {
	if -1 == s.topIndex {
		// 栈空
		return nil, false
	}

	return s.elems[s.topIndex], true
}


// 出栈
func (s *OpStack) PopInt() (int, bool) {
	elem, ok := s.Pop()
	if !ok {
		return 0, ok
	}

	v, ok := elem.(int)
	return v, ok
}

func (s *OpStack) PopReference() (*class.Reference, bool) {
	elem, ok := s.Pop()
	if !ok {
		return nil, ok
	}

	v, ok := elem.(*class.Reference)
	return v, ok
}

func (s *OpStack) GetTopInt() (interface{}, bool) {
	elem, ok := s.GetTop()
	if !ok {
		return 0, ok
	}

	v, ok := elem.(int)
	return v, ok
}

func (s *OpStack) GetTopObject() (*class.Reference, bool) {
	elem, ok := s.GetTop()
	if !ok {
		return nil, ok
	}

	v, ok := elem.(*class.Reference)
	return v, ok
}

func (s *OpStack) GetUntilObject() (*class.Reference, bool) {

	index := s.topIndex
	for {
		if -1 == index {
			break
		}

		v, ok := s.elems[index].(*class.Reference)
		if ok {
			return v, true
		}

		index--
	}

	return nil, false
}
