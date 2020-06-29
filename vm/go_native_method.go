package vm

import (
	"fmt"
	"github.com/wanghongfei/mini-jvm/vm/class"
	"time"
)

func PrintInt(args ...interface{}) interface{} {
	fmt.Println(args[1])

	return true
}

func PrintInt2(args ...interface{}) interface{} {
	fmt.Println(args[1])
	fmt.Println(args[2])

	return true
}

func PrintChar(args ...interface{}) interface{} {
	fmt.Printf("%c\n", args[1])

	return true
}

// 当前协程sleep指定秒数
func ThreadSleep(args ...interface{}) interface{} {
	seconds := args[1].(int)
	time.Sleep(time.Duration(seconds) * time.Second)

	return true
}

// 在新的协程中执行字节码
func ExecuteInThread(args ...interface{}) interface{} {
	// 第一个参数为jvm指针
	jvm := args[0].(*MiniJvm)
	// 第二个参数是实现了Runnalbe接口的对象引用
	objRef := args[1].(*class.Reference)
	// 对象的class定义
	targetClassDef := objRef.Object.DefFile


	// 创建栈帧
	// 把objRef压进去
	opStack := NewOpStack(1)
	opStack.Push(objRef)
	frame := &MethodStackFrame{
		localVariablesTable: nil,
		opStack:             opStack,
		pc:                  0,
	}

	go func() {
		err := jvm.ExecutionEngine.ExecuteWithFrame(targetClassDef, "run", "()V", frame)
		if nil != err {
			fmt.Printf("failed to execute native function 'ExecuteInThread': %v\n", err)
		}
	}()


	return true
}
