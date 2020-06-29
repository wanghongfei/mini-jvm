package vm

import (
	"fmt"
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

