package class

import (
	"fmt"
	"strings"
	"sync"
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

	// 锁
	Monitor sync.Mutex
}


type Object struct {
	// class定义
	DefFile *DefFile

	// 实例数据
	ObjectFields map[string]*ObjectField
}



// 创建对象;
func NewObject(def *DefFile, cl Loader) (*Reference, error) {
	o := new(Object)
	o.DefFile = def

	// 分配字段数据
	o.ObjectFields = make(map[string]*ObjectField)

	// 分配字段, 包括父类里定义的字段
	currentDef := def
	for {
		err := allocateFields(currentDef, o.ObjectFields)
		if nil != err {
			return nil, fmt.Errorf("failed to allcate field for class: %w", err)
		}


		// 父类
		superClassCp := currentDef.ConstPool[currentDef.SuperClass].(*ClassInfoConstInfo)
		superClassFullName := currentDef.ConstPool[superClassCp.FullClassNameIndex].(*Utf8InfoConst).String()
		if "java/lang/Object" == superClassFullName || "java/lang/Exception" == superClassFullName {
			break
		}

		superClassDef, err := cl.LoadClass(superClassFullName)
		if nil != err {
			return nil, fmt.Errorf("failed to load super class '%s' for field allcation: %w", superClassFullName, err)
		}

		currentDef = superClassDef
	}


	//for _, fieldInfo := range def.Fields {
	//	f := new(ObjectField)
	//
	//	// 实例名
	//	name := def.ConstPool[fieldInfo.NameIndex].(*Utf8InfoConst).String()
	//	descriptor := def.ConstPool[fieldInfo.DescriptorIndex].(*Utf8InfoConst).String()
	//	if "I" == descriptor {
	//		f.FieldType = "int"
	//		f.FieldValue = 0
	//
	//	} else {
	//		return nil, fmt.Errorf("unsupported field descriptor '%s'", descriptor)
	//	}
	//
	//	o.ObjectFields[name] = f
	//}

	return &Reference{
		RefType: ReferanceTypeObject,
		Object:  o,
		Array:   nil,
	}, nil
}

func allocateFields(def *DefFile, fields map[string]*ObjectField) error {
	for _, fieldInfo := range def.Fields {
		f := new(ObjectField)

		// 实例名
		name := def.ConstPool[fieldInfo.NameIndex].(*Utf8InfoConst).String()
		descriptor := def.ConstPool[fieldInfo.DescriptorIndex].(*Utf8InfoConst).String()
		if "I" == descriptor {
			f.FieldType = "int"
			f.FieldValue = 0

		} else if "C" == descriptor {
			// char
			f.FieldType = "char"
			f.FieldValue = 'a'

		} else if "[C" == descriptor {
			f.FieldValue = "[]rune"
			f.FieldValue = make([]rune, 0)

		} else if "J" == descriptor {
			f.FieldValue = "long"
			f.FieldValue = 0

		} else if "[Ljava/io/ObjectStreamField;" == descriptor ||
			"Ljava/util/Comparator;" == descriptor {
			// 忽略

		} else {
			return fmt.Errorf("unsupported field descriptor '%s'", descriptor)
		}

		fields[name] = f
	}

	return nil
}

// 创建一个String对象, 用于String字面值常量的创建
func NewStringObject(val []rune, cl Loader) (*Reference, error) {
	stringDef, err := cl.LoadClass("java/lang/String")
	if nil != err {
		return nil, fmt.Errorf("failed to new String object:%w", err)
	}

	obj := &Object{
		DefFile:      stringDef,
		ObjectFields: make(map[string]*ObjectField),
	}

	// 给value和hash这两个最重要的字段赋值
	obj.ObjectFields["value"] = &ObjectField{
		FieldValue: val,
		FieldType:  "[]rune",
	}
	obj.ObjectFields["hash"] = &ObjectField{
		FieldValue: 0,
		FieldType:  "int",
	}

	return &Reference{
		RefType: ReferanceTypeObject,
		Object:  obj,
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

	// 参数列表
	argList := make([]string, 0, 5)

	// 返回类型
	retDesc := descriptor[argDescEndIndex + 1:]

	// 遍历模式
	// 0: 正常模式
	// 1: L模式(解析对象全名, Lxx/xxx/xx;)
	mode := 0
	// sum := 0
	classStartIndex := -1
	for ix, ch := range argDesc {
		// 解析出一个class类型
		if 1 == mode {
			// 处于class解析状态
			if ';' == ch {
				// sum++
				mode = 0

				argList = append(argList, argDesc[classStartIndex:ix])
				classStartIndex = -1
			}

			continue
		}

		if 'L' == ch {
			mode = 1
			classStartIndex = ix
			continue
		}

		argList = append(argList, string(ch))
		// sum++
	}

	return argList, retDesc
}

// 解析方法描述符, 返回方法参数的数量
func ParseArgCount(descriptor string) int {
	argList, _ := ParseMethodDescriptor(descriptor)

	return len(argList)
}

type ObjectField struct {
	// 实例值
	FieldValue interface{}
	FieldType string
}

func NewObjectField(val interface{}) *ObjectField {
	f := new(ObjectField)
	f.FieldValue = val

	switch val.(type) {
	case int:
		f.FieldType = "int"

	case *Reference:
		f.FieldType = "ref"

	case *Array:
		f.FieldType = "arr"

	default:
		f.FieldType = "unknown"
	}

	return f
}

func (f *ObjectField) String() string {
	return fmt.Sprintf("%v", f.FieldValue)
}

type Array struct {
	// 元素类型
	Type byte

	// 数据
	Data []interface{}
}

func NewArray(maxLen int, atype byte) (*Reference, error) {
	//if atype < 4 || atype > 11 {
	//	return nil, fmt.Errorf("unsupported array type '%d'", atype)
	//}

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
