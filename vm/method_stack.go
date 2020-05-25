package vm

import "github.com/wanghongfei/mini-jvm/vm/class"

// 方法栈
type MethodStack struct {
	frames []*MethodStackFrame

	// 永远指向栈顶元素
	topIndex int
}

func NewMethodStack(maxDepth int) *MethodStack {
	return &MethodStack{
		frames:      make([]*MethodStackFrame, maxDepth),
		topIndex:    -1,
	}
}

// 压栈
func (s *MethodStack) Push(frame *MethodStackFrame) bool {
	if s.topIndex == len(s.frames) - 1 {
		// 栈满了
		return false
	}

	s.topIndex++
	s.frames[s.topIndex] = frame

	return true
}

// 出栈
func (s *MethodStack) Pop() (*MethodStackFrame, bool) {
	if -1 == s.topIndex {
		// 栈空
		return nil, false
	}

	data := s.frames[s.topIndex]
	s.topIndex--

	return data, true
}

// 方法栈的栈帧
type MethodStackFrame struct {
	// 本地变量表
	localVariablesTable []interface{}

	// 操作数栈
	opStack *OpStack

	// 程序计数器
	pc int
}

func newMethodStackFrame(opStackDepth int, localVarTableAmount int) *MethodStackFrame {
	return &MethodStackFrame{
		localVariablesTable: make([]interface{}, localVarTableAmount),
		opStack:             NewOpStack(opStackDepth),
		pc:                  0,
	}
}

func (f *MethodStackFrame) GetLocalTableIntAt(index int) int {
	return f.localVariablesTable[index].(int)
}

func (f *MethodStackFrame) GetLocalTableObjectAt(index int) *class.Reference {
	return f.localVariablesTable[index].(*class.Reference)
}
