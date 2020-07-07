package vm

import (
	"fmt"
	"github.com/wanghongfei/mini-jvm/vm/class"
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

func PrintString(args ...interface{}) interface{} {
	strRef := args[1].(*class.Reference)
	field := strRef.Object.ObjectFields["value"]
	strVal := field.FieldValue.([]rune)

	fmt.Printf("%v\n", string(strVal))

	return true
}
