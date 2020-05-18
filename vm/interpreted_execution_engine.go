package vm

import (
	"github.com/wanghongfei/mini-jvm/vm/class"
)

// 解释执行引擎
type InterpretedExecutionEngine struct {
	miniJvm *MiniJvm
}

func (i InterpretedExecutionEngine) Execute(file *class.DefFile) error {
	panic("implement me")
}

func NewInterpretedExecutionEngine(vm *MiniJvm) *InterpretedExecutionEngine {
	return &InterpretedExecutionEngine{miniJvm: vm}
}

