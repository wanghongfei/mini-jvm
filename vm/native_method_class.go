package vm

import (
	"fmt"
	"github.com/wanghongfei/mini-jvm/vm/accflag"
	"github.com/wanghongfei/mini-jvm/vm/class"
	"strings"
)

// Class.getName0()实现
func ClassGetName0(args ...interface{}) interface{} {
	jvm := args[0].(*MiniJvm)
	ref := args[1].(*class.Reference)
	className := ref.Object.DefFile.FullClassName
	className = strings.ReplaceAll(className, "/", ".")

	stringRef, err := class.NewStringObject([]rune(className), jvm.MethodArea)
	if nil != err {
		return fmt.Errorf("failed to create java/lang/String object:%w", err)
	}

	return stringRef
}

// Class.isInterface()
func ClassIsInterface(args ...interface{}) interface{} {
	//  取出class中的accFlag字段
	ref := args[1].(*class.Reference)
	flag := ref.Object.DefFile.AccessFlag
	flagMap := accflag.ParseAccFlags(flag)
	// 判断有没有interface标记位
	if _, ok := flagMap[accflag.Interface]; ok {
		return true
	}

	return false
}

func ClassIsPrimitive(args ...interface{}) interface{} {
	receiver := args[1]
	switch receiver.(type) {
	case int:
		return true
	case int64:
		return true
	case float32:
		return true
	case float64:
		return true

	default:
		return false
	}

}