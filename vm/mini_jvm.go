package vm

import (
	"fmt"
	"github.com/wanghongfei/mini-jvm/vm/class"
	"strings"
)

// VM定义
type MiniJvm struct {
	// 方法区
	MethodArea *MethodArea

	// MainClass全限定性名
	MainClass string

	// 执行引擎
	ExecutionEngine ExecutionEngine
}

type ExecutionEngine interface {
	Execute(file *class.DefFile, methodName string) error
}

func NewMiniJvm(mainClass string, classPaths[] string) (*MiniJvm, error) {
	if "" == mainClass {
		return nil, fmt.Errorf("invalid main class '%s'", mainClass)
	}

	ma, err := NewMethodArea(classPaths)
	if nil != err {
		return nil, fmt.Errorf("unabled to create method area: %w", err)
	}

	vm := &MiniJvm{
		MethodArea: ma,
		MainClass:  strings.ReplaceAll(mainClass, ".", "/"),
	}

	vm.ExecutionEngine = NewInterpretedExecutionEngine(vm)
	return vm, nil
}

// 启动VM
func (m *MiniJvm) Start() error {
	return m.executeMain()
}

// 执行主类
func (m *MiniJvm) executeMain() error {
	mainClassDef, err := m.findDefClass(m.MainClass)
	if nil != err {
		return err
	}

	// 执行
	// log.Printf("main class info: %+v\n", mainClassDef)
	return m.ExecutionEngine.Execute(mainClassDef, "main")
}

func (m *MiniJvm) findDefClass(className string) (*class.DefFile, error) {
	// 从已加载的类中查找
	def, ok := m.MethodArea.ClassMap[className]
	if ok {
		return def, nil
	}

	// 不存在, 触发加载
	def, err := m.MethodArea.LoadClass(className)
	if nil != err {
		return nil, fmt.Errorf("unabled to load class '%s': %w", className, err)
	}

	return def, nil
}
