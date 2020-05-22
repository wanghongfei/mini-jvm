package vm

import (
	"testing"
)

func TestHelloNative(t *testing.T) {
	miniJvm, err := NewMiniJvm("com.fh.HelloNative", []string{"../testclass/"})
	if nil != err {
		t.Fatal(err)
	}

	err = miniJvm.Start()
	if nil != err {
		t.Fatal(err)
	}
}

func TestHelloClass(t *testing.T) {
	miniJvm, err := NewMiniJvm("com.fh.HelloClass", []string{"../testclass/"})
	if nil != err {
		t.Fatal(err)
	}

	err = miniJvm.Start()
	if nil != err {
		t.Fatal(err)
	}

}
