package class

import (
	"fmt"
	"testing"
)

func TestLoadClassFile(t *testing.T) {
	defFile, err := LoadClassFile("../../out/Hello.class")
	if nil != err {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", defFile)
}
