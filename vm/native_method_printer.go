package vm

import (
	"fmt"
	"github.com/wanghongfei/mini-jvm/utils"
	"github.com/wanghongfei/mini-jvm/vm/class"
)

func PrintInt(args ...interface{}) interface{} {
	fmt.Println(args[2])

	return nil
}

func PrintInt2(args ...interface{}) interface{} {
	fmt.Println(args[2])
	fmt.Println(args[3])

	return nil
}

func PrintChar(args ...interface{}) interface{} {
	fmt.Printf("%c\n", args[2])

	return nil
}

func PrintString(args ...interface{}) interface{} {
	strRef := args[2].(*class.Reference)
	field := strRef.Object.ObjectFields["value"]
	strArrayRef := field.FieldValue.(*class.Reference)

	runeArr := utils.InterfaceArrayToRuneArray(strArrayRef.Array.Data)

	fmt.Printf("%v\n", string(runeArr))

	return nil
}

func PrintBoolean(args ...interface{}) interface{} {
	boolInt := args[2].(int)
	if 0 == boolInt {
		fmt.Println("false")

	} else {
		fmt.Println("true")
	}

	return nil
}
