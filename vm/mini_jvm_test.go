package vm

import (
	"testing"
)

func TestHelloNative(t *testing.T) {
	miniJvm, err := NewMiniJvm("com.fh.ForLoopPrintTest", []string{"../testclass/"})
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
	miniJvm, err := NewMiniJvm("com.fh.NewSimpleObjectTest", []string{"../testclass/"})
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
	miniJvm, err := NewMiniJvm("com.fh.MethodReloadTest", []string{"../testclass/"})
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
	miniJvm, err := NewMiniJvm("com.fh.ClassExtendTest", []string{"../testclass/"})
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
	miniJvm, err := NewMiniJvm("com.fh.RecursionTest", []string{"../testclass/"})
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
	miniJvm, err := NewMiniJvm("com.fh.IfTest", []string{"../testclass/"})
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
	if -301 != miniJvm.DebugPrintHistory[99] {
		t.FailNow()
	}
}

func TestIntArray(t *testing.T) {
	miniJvm, err := NewMiniJvm("com.fh.ArrayTest", []string{"../testclass/"})
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
}

func TestInterface(t *testing.T) {
	miniJvm, err := NewMiniJvm("com.fh.InterfaceTest", []string{"../testclass/"})
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
