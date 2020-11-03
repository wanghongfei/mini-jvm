package vm

import (
	"fmt"
	"github.com/wanghongfei/mini-jvm/utils"
	"github.com/wanghongfei/mini-jvm/vm/class"
	"time"
)

const (
	THREAD_STATUS_CREATED = 0
	THREAD_STATUS_RUNNING = 0
	THREAD_STATUS_FINISHED = 0
)

// java线程对应go里的表示
type MiniThread struct {
	Jvm *MiniJvm
	JavaObjRef *class.Reference

	// 线程状态
	// 0: created
	// 1: running
	// 2: finished
	Status int
}

func (t *MiniThread) Start() {
	// 创建栈帧
	// 把objRef压进去
	opStack := NewOpStack(1)
	opStack.Push(t.JavaObjRef)
	frame := &MethodStackFrame{
		localVariablesTable: nil,
		opStack:             opStack,
		pc:                  0,
	}

	go func() {
		t.Status = THREAD_STATUS_RUNNING

		// 防止进程崩溃
		defer func() {
			r := recover()
			if nil != r {
				fmt.Printf("goroutine recovered: %v\n", r)
			}
		}()

		defer func() {
			t.Status = THREAD_STATUS_FINISHED
		}()

		err := t.Jvm.ExecutionEngine.ExecuteWithFrame(t.JavaObjRef.Object.DefFile, "run", "()V", frame, false)
		if nil != err {
			if expRef, ok := err.(*ExceptionThrownError); ok {
				// 底层抛出了没有捕获的异常
				utils.LogErrorPrintf("thread exit due to thrown exception: %v\n", expRef.ExceptionRef.Object.DefFile.FullClassName)
				return
			}

			utils.LogInfoPrintf("failed to execute native function 'ExecuteInThread': %v\n", err)
		}
	}()
}

// 当前协程sleep指定秒数
func ThreadSleep(args ...interface{}) interface{} {
	seconds := args[2].(int)
	time.Sleep(time.Duration(seconds) * time.Second)

	return nil
}

// 在新的协程中执行字节码
func ExecuteInThread(args ...interface{}) interface{} {
	// 第一个参数为jvm指针
	jvm := args[0].(*MiniJvm)
	// 第三个参数是实现了Runnalbe接口的对象引用
	objRef := args[2].(*class.Reference)

	miniThread := &MiniThread{
		Jvm:        jvm,
		JavaObjRef: objRef,
		Status: THREAD_STATUS_CREATED,
	}
	miniThread.Start()

	return nil
}
