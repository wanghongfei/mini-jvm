package class

import (
	"fmt"
	"strings"
)

const (
	ReferanceTypeObject = byte(0)
	ReferanceTypeArray = byte(1)
)

// 表达Java中的引用类型
type Reference struct {
	// 引用的类型
	// 0: object
	// 1: array
	RefType byte

	Object *Object
	Array *Array
}

type Array struct {
	// 元素类型
	Type byte

	// 数据
	Data []interface{}
}

type Object struct {
	// class定义
	DefFile *DefFile

	// 实例数据
	ObjectFields map[string]*ObjectField
}


type ObjectField struct {
	// 实例值
	FieldValue interface{}
	FieldType string
}

// 创建对象;
func NewObject(def *DefFile) (*Reference, error) {
	o := new(Object)
	o.DefFile = def

	// 分配字段数据
	o.ObjectFields = make(map[string]*ObjectField)

	// 初始化字段
	for _, fieldInfo := range def.Fields {
		f := new(ObjectField)

		// 实例名
		name := def.ConstPool[fieldInfo.NameIndex].(*Utf8InfoConst).String()
		descriptor := def.ConstPool[fieldInfo.DescriptorIndex].(*Utf8InfoConst).String()
		if "I" == descriptor {
			f.FieldType = "int"
			f.FieldValue = 0

		} else {
			return nil, fmt.Errorf("unsupported field descriptor '%s'", descriptor)
		}

		o.ObjectFields[name] = f
	}

	return &Reference{
		RefType: ReferanceTypeObject,
		Object:  o,
		Array:   nil,
	}, nil
}


// 解析方法描述符;
// ret1: 参数列表
// ret2: 返回类型
func ParseMethodDescriptor(descriptor string) ([]string, string) {
	// 提取参数列表
	argDescEndIndex := strings.Index(descriptor, ")")
	argDesc := descriptor[1:argDescEndIndex]

	// 解析参数列表
	argList := make([]string, 0, 5)
	for _, ch := range argDesc {
		argList = append(argList, string(ch))
	}

	retDesc := descriptor[argDescEndIndex + 1:]

	return argList, retDesc
}


func NewArray(maxLen int, atype byte) (*Reference, error) {
	if atype < 4 || atype > 11 {
		return nil, fmt.Errorf("unsupported array type '%d'", atype)
	}

	arr := &Array{
		Type: atype,
		Data: make([]interface{}, maxLen),
	}

	return &Reference{
		RefType: ReferanceTypeArray,
		Object:  nil,
		Array:   arr,
	}, nil
}
