package class

import (
	"github.com/wanghongfei/mini-jvm/utils"
	"io"
)

type ClassInfoConstInfo struct {
	Tag uint8
	FullClassNameIndex uint16
}

func ReadClassInfoConst(reader io.Reader, tag uint8) (*ClassInfoConstInfo, error) {
	cpInfo := new(ClassInfoConstInfo)
	cpInfo.Tag = tag

	classIndex, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, err
	}
	cpInfo.FullClassNameIndex = classIndex

	return cpInfo, nil
}

// 字段引用常量
type FieldRefConstInfo struct {
	Tag uint8
	ClassIndex uint16
	NameAndTypeIndex uint16
}

func ReadFieldRefConst(reader io.Reader, tag uint8) (*FieldRefConstInfo, error) {
	cpInfo := new(FieldRefConstInfo)
	cpInfo.Tag = tag

	classIndex, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, err
	}
	cpInfo.ClassIndex = classIndex

	typeIndex, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, err
	}
	cpInfo.NameAndTypeIndex = typeIndex

	return cpInfo, nil
}

type IntegerInfoConst struct {
	Tag uint8
	Bytes uint32
}


func ReadIntegerInfoConst(reader io.Reader, tag uint8) (*IntegerInfoConst, error) {
	cpInfo := new(IntegerInfoConst)
	cpInfo.Tag = tag

	b, err := utils.ReadInt32(reader)
	if nil != err {
		return nil, err
	}
	cpInfo.Bytes = b

	return cpInfo, nil
}

// 方法引用常量
type MethodRefConstInfo struct {
	Tag uint8
	ClassIndex uint16
	NameAndTypeIndex uint16
}

func ReadMethodRefConst(reader io.Reader, tag uint8) (*MethodRefConstInfo, error) {
	cpInfo := new(MethodRefConstInfo)
	cpInfo.Tag = tag

	classIndex, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, err
	}
	cpInfo.ClassIndex = classIndex

	typeIndex, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, err
	}
	cpInfo.NameAndTypeIndex = typeIndex

	return cpInfo, nil
}

type NameAndTypeConst struct {
	Tag uint8
	NameIndex uint16
	DescIndex uint16
}

func ReadNameAndTypeConst(reader io.Reader, tag uint8) (*NameAndTypeConst, error) {
	cpInfo := new(NameAndTypeConst)
	cpInfo.Tag = tag

	name, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, err
	}
	cpInfo.NameIndex = name

	desc, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, err
	}
	cpInfo.DescIndex = desc

	return cpInfo, nil
}

type StringInfoConst struct {
	Tag uint8
	StringIndex uint16
}

func ReadStringInfoConst(reader io.Reader, tag uint8) (*StringInfoConst, error) {
	cpInfo := new(StringInfoConst)
	cpInfo.Tag = tag

	index, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, err
	}
	cpInfo.StringIndex = index

	return cpInfo, nil
}

type Utf8InfoConst struct {
	Tag uint8
	Length uint16
	Bytes []byte
}

func ReadUtf8InfoConst(reader io.Reader, tag uint8) (*Utf8InfoConst, error) {
	cpInfo := new(Utf8InfoConst)
	cpInfo.Tag = tag

	length, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, err
	}
	cpInfo.Length = length

	utf8Buf := make([]byte, length)
	_, err = reader.Read(utf8Buf)
	if nil != err {
		return nil, err
	}
	cpInfo.Bytes = utf8Buf

	return cpInfo, nil
}

func (o *Utf8InfoConst) String() string {
	return string(o.Bytes)
}

type InterfaceMethodConst struct {
	Tag uint8
	InterfaceClassIndex uint16
	NameAndTypeIndex uint16
}

func ReadInterfaceMethodConst(reader io.Reader, tag uint8) (*InterfaceMethodConst, error) {
	cpInfo := new(InterfaceMethodConst)
	cpInfo.Tag = tag

	interfaceIndex, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, err
	}
	cpInfo.InterfaceClassIndex = interfaceIndex

	nameAndTypeIndex, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, err
	}
	cpInfo.NameAndTypeIndex = nameAndTypeIndex

	return cpInfo, nil
}