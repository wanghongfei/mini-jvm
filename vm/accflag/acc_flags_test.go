package accflag

import (
	"fmt"
	"testing"
)

func TestParseAccFlags(t *testing.T) {
	bits := 0x0001
	fmt.Println(ParseAccFlags(uint16(bits)))
}
