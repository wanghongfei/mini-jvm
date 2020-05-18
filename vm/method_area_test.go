package vm

import (
	"fmt"
	"testing"
)

func TestMethodArea_LoadClass(t *testing.T) {
	ma, err := NewMethodArea([]string{"../out"})
	if nil != err {
		t.Fatal(err)
	}
	_, err = ma.LoadClass("Hello")
	if nil != err {
		t.Fatal(err)
	}

	_, err = ma.LoadClass("com/fh/vo/Student")
	if nil != err {
		t.Fatal(err)
	}

	fmt.Printf("%v\n", ma)
}
