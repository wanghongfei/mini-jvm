package vm

import (
	"github.com/wanghongfei/mini-jvm/vm/class"
	"testing"
)

// 改成自己电脑中rt.jar的路径
var rtJarPath = "/Library/Java/JavaVirtualMachines/jdk1.8.0_181.jdk/Contents/Home/jre/lib/rt.jar"

func TestHelloNative(t *testing.T) {
	miniJvm, err := NewMiniJvm("com.fh.ForLoopPrintTest", []string{"../testcase/classes", "../mini-lib/classes", "/Library/Java/JavaVirtualMachines/jdk1.8.0_181.jdk/Contents/Home/jre/lib/rt.jar"})
	if nil != err {
		t.Fatal(err)
	}

	err = miniJvm.Start()
	if nil != err {
		t.Fatal(err)
	}

	// assert
	if 5050 != miniJvm.DebugPrintHistory[0] {
		t.FailNow()
	}
}

func TestHelloClass(t *testing.T) {
	miniJvm, err := NewMiniJvm("com.fh.NewSimpleObjectTest", []string{"../testcase/classes", "../mini-lib/classes", rtJarPath})
	if nil != err {
		t.Fatal(err)
	}

	err = miniJvm.Start()
	if nil != err {
		t.Fatal(err)
	}

	// assert
	if 5050 != miniJvm.DebugPrintHistory[0] {
		t.FailNow()
	}
}

func TestHelloMethod(t *testing.T) {
	miniJvm, err := NewMiniJvm("com.fh.MethodReloadTest", []string{"../testcase/classes", "../mini-lib/classes", rtJarPath})
	if nil != err {
		t.Fatal(err)
	}

	err = miniJvm.Start()
	if nil != err {
		t.Fatal(err)
	}

	// assert
	if 300 != miniJvm.DebugPrintHistory[0] {
		t.FailNow()
	}
}

func TestClassExtend(t *testing.T) {
	miniJvm, err := NewMiniJvm("com.fh.ClassExtendTest", []string{"../testcase/classes", "../mini-lib/classes", rtJarPath})
	if nil != err {
		t.Fatal(err)
	}

	err = miniJvm.Start()
	if nil != err {
		t.Fatal(err)
	}

	// assert
	if 1 != miniJvm.DebugPrintHistory[0] {
		t.FailNow()
	}
	if 10 != miniJvm.DebugPrintHistory[1] {
		t.FailNow()
	}
}

func TestRecursion(t *testing.T) {
	miniJvm, err := NewMiniJvm("com.fh.RecursionTest", []string{"../testcase/classes", "../mini-lib/classes", rtJarPath})
	if nil != err {
		t.Fatal(err)
	}

	err = miniJvm.Start()
	if nil != err {
		t.Fatal(err)
	}

	// assert
	if 1 != miniJvm.DebugPrintHistory[0] {
		t.FailNow()
	}
	if 100 != miniJvm.DebugPrintHistory[99] {
		t.FailNow()
	}
}

func TestIfTest(t *testing.T) {
	miniJvm, err := NewMiniJvm("com.fh.IfTest", []string{"../testcase/classes", "../mini-lib/classes", rtJarPath})
	if nil != err {
		t.Fatal(err)
	}

	err = miniJvm.Start()
	if nil != err {
		t.Fatal(err)
	}

	// assert
	if -301 != miniJvm.DebugPrintHistory[0] {
		t.FailNow()
	}
	if -301 != miniJvm.DebugPrintHistory[1] {
		t.FailNow()
	}
}

func TestArray(t *testing.T) {
	miniJvm, err := NewMiniJvm("com.fh.ArrayTest", []string{"../testcase/classes", "../mini-lib/classes", rtJarPath})
	if nil != err {
		t.Fatal(err)
	}

	err = miniJvm.Start()
	if nil != err {
		t.Fatal(err)
	}


	// assert
	if 1 != miniJvm.DebugPrintHistory[0] {
		t.FailNow()
	}
	if 2 != miniJvm.DebugPrintHistory[1] {
		t.FailNow()
	}
	if 3 != miniJvm.DebugPrintHistory[2] {
		t.FailNow()
	}
	if 4 != miniJvm.DebugPrintHistory[3] {
		t.FailNow()
	}
	if int('好') != miniJvm.DebugPrintHistory[4] {
		t.FailNow()
	}
	if int('吗') != miniJvm.DebugPrintHistory[5] {
		t.FailNow()
	}

	if 3 != miniJvm.DebugPrintHistory[6] {
		t.FailNow()
	}
	if 5 != miniJvm.DebugPrintHistory[7] {
		t.FailNow()
	}
}

func TestObjectArray(t *testing.T) {
	miniJvm, err := NewMiniJvm("com.fh.ObjectArrayTest", []string{"../testcase/classes", "../mini-lib/classes", rtJarPath})
	if nil != err {
		t.Fatal(err)
	}

	err = miniJvm.Start()
	if nil != err {
		t.Fatal(err)
	}

	// assert
	if 0 != miniJvm.DebugPrintHistory[0] {
		t.FailNow()
	}
	if 100 != miniJvm.DebugPrintHistory[1] {
		t.FailNow()
	}
}

func TestInterface(t *testing.T) {
	miniJvm, err := NewMiniJvm("com.fh.InterfaceTest", []string{"../testcase/classes", "../mini-lib/classes", rtJarPath})
	if nil != err {
		t.Fatal(err)
	}

	err = miniJvm.Start()
	if nil != err {
		t.Fatal(err)
	}

	// assert
	if 100 != miniJvm.DebugPrintHistory[0] {
		t.FailNow()
	}
	if 100 != miniJvm.DebugPrintHistory[1] {
		t.FailNow()
	}
	if 500 != miniJvm.DebugPrintHistory[2] {
		t.FailNow()
	}
}

func TestException(t *testing.T) {
	miniJvm, err := NewMiniJvm("com.fh.ExceptionTest", []string{"../testcase/classes", "../mini-lib/classes", rtJarPath})
	if nil != err {
		t.Fatal(err)
	}

	err = miniJvm.Start()
	if nil != err {
		t.Fatal(err)
	}

	// assert
	if 10 != miniJvm.DebugPrintHistory[0] {
		t.FailNow()
	}
	if 20 != miniJvm.DebugPrintHistory[1] {
		t.FailNow()
	}
	if 30 != miniJvm.DebugPrintHistory[2] {
		t.FailNow()
	}
}


func TestHanoi(t *testing.T) {
	miniJvm, err := NewMiniJvm("com.fh.Hanoi", []string{"../testcase/classes", "../mini-lib/classes", rtJarPath})
	if nil != err {
		t.Fatal(err)
	}

	err = miniJvm.Start()
	if nil != err {
		t.Fatal(err)
	}

	if 127 != miniJvm.DebugPrintHistory[0] {
		t.FailNow()
	}
}

func TestStaticField(t *testing.T) {
	miniJvm, err := NewMiniJvm("com.fh.StaticFieldTest", []string{"../testcase/classes", "../mini-lib/classes", rtJarPath})
	if nil != err {
		t.Fatal(err)
	}

	err = miniJvm.Start()
	if nil != err {
		t.Fatal(err)
	}

	// assert
	if 100 != miniJvm.DebugPrintHistory[0].(*class.ObjectField).FieldValue {
		t.FailNow()
	}
	if 400 != miniJvm.DebugPrintHistory[1].(*class.ObjectField).FieldValue {
		t.FailNow()
	}
}

func TestThread(t *testing.T) {
	miniJvm, err := NewMiniJvm("com.fh.thread.ThreadTest", []string{"../testcase/classes", "../mini-lib/classes", rtJarPath})
	if nil != err {
		t.Fatal(err)
	}

	err = miniJvm.Start()
	if nil != err {
		t.Fatal(err)
	}

}

func TestString(t *testing.T) {
	miniJvm, err := NewMiniJvm("com.fh.StringTest", []string{"../testcase/classes", "../mini-lib/classes", rtJarPath})
	if nil != err {
		t.Fatal(err)
	}

	err = miniJvm.Start()
	if nil != err {
		t.Fatal(err)
	}

	// assert
	val := miniJvm.DebugPrintHistory[0].(*class.Reference).Object.ObjectFields["value"].FieldValue.([]rune)
	if "hello, 世界" != string(val) {
		t.FailNow()
	}

}

func TestObjectLoading(t *testing.T) {
	miniJvm, err := NewMiniJvm("com.fh.ObjectTest", []string{"../testcase/classes", "../mini-lib/classes", rtJarPath})
	if nil != err {
		t.Fatal(err)
	}

	err = miniJvm.Start()
	if nil != err {
		t.Fatal(err)
	}
}

func TestReflection(t *testing.T) {
	miniJvm, err := NewMiniJvm("com.fh.ReflectionTest", []string{"../testcase/classes", "../mini-lib/classes", rtJarPath})
	if nil != err {
		t.Fatal(err)
	}

	err = miniJvm.Start()
	if nil != err {
		t.Fatal(err)
	}

}