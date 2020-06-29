package vm

import (
	"strings"
)

// JVM的本地方法, 即go函数;
// 参数args[0]固定为MiniJVM的指针
type NativeFunction func(args ...interface{}) interface{}

type NativeMethodInfo struct {
	// 方法名
	Name string

	// 描述符;
	// String getRealnameByIdAndNickname(int id,String name) 的描述符为 (ILjava/lang/String;)Ljava/lang/String;
	Descriptor string

	// 对应的go函数
	EntryFunc NativeFunction
}

// 解析方法描述符, 返回方法参数的数量
// 注意, 不支持方法描述符中含有对象类型
func (info *NativeMethodInfo) ParseArgCount() int {
	argAndRetDesciptor := strings.Split(info.Descriptor, ")")
	argDescriptor := argAndRetDesciptor[0][1:]
	//
	//return len(argDescriptor)


	// 遍历模式
	// 0: 正常模式
	// 1: L模式(解析对象全名, Lxx/xxx/xx;)
	mode := 0
	sum := 0
	for _, ch := range argDescriptor {
		// 解析出一个class类型
		if 1 == mode {
			// 处于class解析状态
			if ';' == ch {
				sum++
				mode = 0
			}

			continue
		}

		if 'L' == ch {
			mode = 1
			continue
		}

		sum++
	}

	return sum

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
func (t *NativeMethodTable) RegisterMethod(methodName string, descriptor string, goFunc NativeFunction) {
	key := t.genKey(methodName, descriptor)
	t.MethodInfoMap[key] = &NativeMethodInfo{
		Name:       methodName,
		Descriptor: descriptor,
		EntryFunc:  goFunc,
	}
}

// 查本地方法表, 找出目标go函数
func (t *NativeMethodTable) FindMethod(name string, descriptor string) (NativeFunction, int) {
	f, ok := t.MethodInfoMap[t.genKey(name, descriptor)]
	if !ok {
		return nil, -1
	}

	return f.EntryFunc, f.ParseArgCount()
}


func (t *NativeMethodTable) genKey(name string, descriptor string) string {
	return name + "=>" + descriptor
}
