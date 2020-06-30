package vm

import (
	"github.com/wanghongfei/mini-jvm/vm/class"
	"strings"
)

// JVM的本地方法, 即go函数;
// 参数args[0]固定为MiniJVM的指针
type NativeFunction func(args ...interface{}) interface{}

type NativeMethodInfo struct {
	// 方法名
	Name string

	// 类的全名
	FullClassName string

	// 描述符;
	// String getRealnameByIdAndNickname(int id,String name) 的描述符为 (ILjava/lang/String;)Ljava/lang/String;
	Descriptor string

	// 对应的go函数
	EntryFunc NativeFunction
}


// 本地方法表
type NativeMethodTable struct {
	MethodInfoMap map[string]*NativeMethodInfo
}

func NewNativeMethodTable() *NativeMethodTable {
	return &NativeMethodTable{MethodInfoMap: map[string]*NativeMethodInfo{}}
}

// 注册本地方法
// methodName: 方法名
// descriptor: 方法在JVM中的描述符
func (t *NativeMethodTable) RegisterMethod(className string, methodName string, descriptor string, goFunc NativeFunction) {
	key := t.genKey(strings.ReplaceAll(className, ".", "/"), methodName, descriptor)
	t.MethodInfoMap[key] = &NativeMethodInfo{
		Name:       methodName,
		Descriptor: descriptor,
		EntryFunc:  goFunc,
	}
}

// 查本地方法表, 找出目标go函数
func (t *NativeMethodTable) FindMethod(className, name string, descriptor string) (NativeFunction, int) {
	key := t.genKey(className, name, descriptor)
	f, ok := t.MethodInfoMap[key]
	if !ok {
		return nil, -1
	}

	return f.EntryFunc, class.ParseArgCount(f.Descriptor)
}


func (t *NativeMethodTable) genKey(className, methodName string, descriptor string) string {
	return className + ";" + methodName + ";" + descriptor
}
