package vm

import (
	"testing"
)

func TestHelloNative(t *testing.T) {
	miniJvm, err := NewMiniJvm("HelloNative", []string{"../testclass/"})
	if nil != err {
		t.Fatal(err)
	}

	err = miniJvm.Start()
	if nil != err {
		t.Fatal(err)
	}
}
