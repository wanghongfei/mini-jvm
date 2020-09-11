package vm

import "github.com/wanghongfei/mini-jvm/vm/class"

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
	elem := f.localVariablesTable[index]
	if nil == elem {
		return nil
	}

	return elem.(*class.Reference)
}
