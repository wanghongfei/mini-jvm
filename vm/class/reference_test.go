package class

import (
	"fmt"
	"testing"
)

func TestParseMethodDescriptor(t *testing.T) {
	args, ret := ParseMethodDescriptor("([CI)[C")
	fmt.Println(args)
	fmt.Println(ret)
	if args[0] != "[C" || args[1] != "I" {
		t.FailNow()
	}

	args, ret = ParseMethodDescriptor("(II[CI)V")
	fmt.Println(args)
	fmt.Println(ret)
	if args[0] != "I" || args[1] != "I" || args[2] != "[C" || args[3] != "I" {
		t.FailNow()
	}
}
