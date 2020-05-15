package class

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/wanghongfei/mini-jvm/utils"
	"io"
)

// class文件定义
type DefFile struct {
	MagicNumber uint32

	MinorVersion uint16
	MajorVersion uint16

	// 常量池数量
	ConstPoolCount uint16
	// 常量池
	ConstPool []interface{}

	// 访问标记
	AccessFlag uint16
	// 当前类在常量池的索引
	ThisClass uint16
	// 父类索引
	SuperClass uint16

	// 接口
	InterfacesCount uint16
	Interfaces []uint16

	// 字段
	FieldsCount uint16
	Fields []*FieldInfo

	// 方法
	MethodCount uint16
	Methods []*MethodInfo

	// 属性
	AttrCount uint16
	Attrs []interface{}
}




const JVM_CLASS_FILE_MAGIC_NUMBER = 0xCAFEBABE

// 加载class文件
func LoadClassFile(classPath string) (*DefFile, error) {
	classBuf, err := utils.ReadAllFromFile(classPath)
	if nil != err {
		return nil, fmt.Errorf("failed to read class file, %w", err)
	}

	defFile := new(DefFile)

	bufReader := bytes.NewReader(classBuf)

	// 魔术数
	defFile.MagicNumber, err = utils.ReadInt32(bufReader)
	if nil != err {
		return nil, fmt.Errorf("failed to load magic number, %w", err)
	}
	if defFile.MagicNumber != JVM_CLASS_FILE_MAGIC_NUMBER {
		return nil, errors.New("not a JVM class file")
	}

	// 副版本号
	defFile.MinorVersion, err = utils.ReadInt16(bufReader)
	if nil != err {
		return nil, fmt.Errorf("failed to load minor version, %w", err)
	}
	// 主版本号
	defFile.MajorVersion, err = utils.ReadInt16(bufReader)
	if nil != err {
		return nil, fmt.Errorf("failed to load minor version, %w", err)
	}

	// 常量池数量
	defFile.ConstPoolCount, err = utils.ReadInt16(bufReader)
	if nil != err {
		return nil, fmt.Errorf("failed to load const pool count, %w", err)
	}

	// 常量池
	defFile.ConstPool, err = readConstPool(bufReader, int(defFile.ConstPoolCount) - 1)
	if nil != err {
		return nil, fmt.Errorf("failed to load const pool, %w", err)
	}

	// 访问标记
	defFile.AccessFlag, err = utils.ReadInt16(bufReader)
	if nil != err {
		return nil, fmt.Errorf("failed to load access_flag, %w", err)
	}

	// 当期类在常量池的索引
	defFile.ThisClass, err = utils.ReadInt16(bufReader)
	if nil != err {
		return nil, fmt.Errorf("failed to load this_class, %w", err)
	}

	// 父类索引
	defFile.SuperClass, err = utils.ReadInt16(bufReader)
	if nil != err {
		return nil, fmt.Errorf("failed to load super_class, %w", err)
	}

	// 接口数量
	defFile.InterfacesCount, err = utils.ReadInt16(bufReader)
	if nil != err {
		return nil, fmt.Errorf("failed to load interfaces_count, %w", err)
	}

	// 接口索引
	defFile.Interfaces = make([]uint16, 0, defFile.InterfacesCount)
	for ix := 0; ix < int(defFile.InterfacesCount); ix++ {
		index, err := utils.ReadInt16(bufReader)
		if nil != err {
			return nil, fmt.Errorf("failed to load interfaces_index, %w", err)
		}

		defFile.Interfaces = append(defFile.Interfaces, index)
	}

	// 字段表长度
	defFile.FieldsCount, err = utils.ReadInt16(bufReader)
	if nil != err {
		return nil, fmt.Errorf("failed to load fields_count, %w", err)
	}

	// 读取字段表
	defFile.Fields = make([]*FieldInfo, 0, defFile.FieldsCount)
	for ix := 0; ix < int(defFile.FieldsCount); ix++ {
		f, err := defFile.ReadFieldInfo(bufReader)
		if nil != err {
			return nil, fmt.Errorf("failed to load field_info, %w", err)
		}

		defFile.Fields = append(defFile.Fields, f)
	}


	// 方法表长度
	defFile.MethodCount, err = utils.ReadInt16(bufReader)
	if nil != err {
		return nil, fmt.Errorf("failed to load method_count: %w", err)
	}

	// 读取方法表
	defFile.Methods = make([]*MethodInfo, 0, defFile.MethodCount)
	for ix := 0; ix < int(defFile.MethodCount); ix++ {
		m, err := defFile.ReadMethodInfo(bufReader)
		if nil != err {
			return nil, fmt.Errorf("failed to load method_info: %w", err)
		}

		defFile.Methods = append(defFile.Methods, m)
	}

	// 属性长度
	attrCount, err := utils.ReadInt16(bufReader)
	if nil != err {
		return nil, err
	}
	defFile.AttrCount = attrCount

	// 读取属性表
	defFile.Attrs = make([]interface{}, 0, attrCount)
	for ix := 0; ix < int(attrCount); ix++ {
		attr, err := defFile.ReadAttr(bufReader)
		if nil != err {
			return nil, fmt.Errorf("failed to read attr, %w", err)
		}

		defFile.Attrs = append(defFile.Attrs, attr)
	}

	return defFile, nil
}

// 解析常量池
func readConstPool(bufReader io.Reader, cpCount int) ([]interface{}, error) {
	cpInfos := make([]interface{}, 1, cpCount)

	for ix := 0; ix < cpCount; ix++ {
		// 读取一个常量

		// 读取tag
		tag, err := utils.ReadInt8(bufReader)
		if nil != err {
			return nil, err
		}

		// 根据tag判断常量类型
		switch tag {
		case 1:
			info, err := ReadUtf8InfoConst(bufReader, tag)
			if nil != err {
				return nil, err
			}
			cpInfos = append(cpInfos, info)

		case 3:
			info, err := ReadIntegerInfoConst(bufReader, tag)
			if nil != err {
				return nil, err
			}
			cpInfos = append(cpInfos, info)

		case 7:
			info, err := ReadClassInfoConst(bufReader, tag)
			if nil != err {
				return nil, err
			}
			cpInfos = append(cpInfos, info)

		case 9:
			info, err := ReadFieldRefConst(bufReader, tag)
			if nil != err {
				return nil, err
			}
			cpInfos = append(cpInfos, info)

		case 10:
			info, err := ReadMethodRefConst(bufReader, tag)
			if nil != err {
				return nil, err
			}
			cpInfos = append(cpInfos, info)

		case 12:
			info, err := ReadNameAndTypeConst(bufReader, tag)
			if nil != err {
				return nil, err
			}
			cpInfos = append(cpInfos, info)


		default:
			return nil, fmt.Errorf("invalid cp tag %d", tag)
		}
	}

	return cpInfos, nil
}


