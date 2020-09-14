package class

import (
	"errors"
	"fmt"
	"github.com/wanghongfei/mini-jvm/utils"
	"io"
)

func (c *DefFile) ReadAttr(reader io.Reader) (interface{}, error) {
	// 读取属性名index
	nameIndex, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, fmt.Errorf("failed to read attr_name_index, %w", err)
	}

	// 从常量池中查找属性名
	attrNameObject, err := c.GetFromConstPool(int(nameIndex))
	if nil != err {
		return nil, fmt.Errorf("failed to get const from const pool, %w", err)
	}

	utf8Const, ok := attrNameObject.(*Utf8InfoConst)
	if !ok {
		return nil, fmt.Errorf("const pool at index %d is not an utf8 const", nameIndex)
	}

	// 取出属性名
	attrName := utf8Const.String()
	if "Code" == attrName {
		codeAttr, err := c.ReadCodeAttr(reader)
		if nil != err {
			return nil, fmt.Errorf("failed to read code attr, %w", err)
		}

		return codeAttr, nil

	} else if "ConstantValue" == attrName {
		valAttr, err := ReadConstantValueAttr(reader)
		if nil != err {
			return nil, fmt.Errorf("failed to read constant_value attr: %w", err)
		}

		return valAttr, nil

	} else if "LineNumberTable" == attrName {
		lineAttr, err := ReadLineNumberAttr(reader)
		if nil != err {
			return nil, fmt.Errorf("failed to read LineNumberTable attr: %w", err)
		}

		return lineAttr, nil

	} else if "SourceFile" == attrName {
		srcAttr, err := ReadSourceFileAttr(reader)
		if nil != err {
			return nil, fmt.Errorf("failed to read SourceFile attr: %w", err)
		}

		return srcAttr, nil

	} else if "StackMapTable" == attrName ||
		"Signature" == attrName ||
		"Deprecated" == attrName ||
		"RuntimeVisibleAnnotations" == attrName ||
		"Exceptions" == attrName ||
		"BootstrapMethods" == attrName {
		// 跳过此属性
		err := c.skipAttr(reader)
		if nil != err {
			return nil, fmt.Errorf("failed to skip StackMapTable attr: %w", err)
		}

		return struct{}{}, nil

	} else if "InnerClasses" == attrName {
		innerAttr, err := ReadInnerClassAttr(reader)
		if nil != err {
			return nil, fmt.Errorf("failed to read InnerClasses attr: %w", err)
		}

		return innerAttr, nil
	}

	return nil, fmt.Errorf("unsupported attr type '%s'", attrName)
}

// 从常量池中取出常量
func (c *DefFile) GetFromConstPool(index int) (interface{}, error) {
	if index >= len(c.ConstPool) {
		return nil, errors.New("cp index out of bound")
	}

	return c.ConstPool[index], nil
}

func (c *DefFile) skipAttr(reader io.Reader) error {
	attrLen, err := utils.ReadInt32(reader)
	if nil != err {
		return fmt.Errorf("failed to read skip len: %w", err)
	}

	for ix := 0; ix < int(attrLen); ix++ {
		utils.ReadInt8(reader)
	}

	return nil
}

// code属性
type CodeAttr struct {
	AttrLength uint32

	MaxStack uint16
	MaxLocals uint16

	// 字节码长度
	CodeLength uint32
	Code []byte

	// 异常表
	ExceptionTableLength uint16
	ExceptionTable []*ExceptionTable

	AttrCount uint16
	Attrs []interface{}
}

func (c *CodeAttr) String() string {
	return "Code"
}

type ExceptionTable struct {
	StartPc uint16
	EndPc uint16
	HandlerPc uint16
	CatchType uint16
}

func (c *DefFile) ReadCodeAttr(reader io.Reader) (*CodeAttr, error) {
	attr := new(CodeAttr)

	length, err := utils.ReadInt32(reader)
	if nil != err {
		return nil, err
	}
	attr.AttrLength = length

	maxStack, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, err
	}
	attr.MaxStack = maxStack

	maxLocals, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, err
	}
	attr.MaxLocals = maxLocals

	// 字节码长度
	codeLen, err := utils.ReadInt32(reader)
	if nil != err {
		return nil, err
	}
	attr.CodeLength = codeLen

	// 字节码
	codeBuf := make([]byte, codeLen)
	_, err = reader.Read(codeBuf)
	if nil != err {
		return nil, err
	}
	attr.Code = codeBuf

	// 异常表长度
	expTableLen, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, err
	}
	attr.ExceptionTableLength = expTableLen

	// 异常表
	expLen := int(expTableLen)
	attr.ExceptionTable = make([]*ExceptionTable, 0, expLen)
	for ix := 0; ix < expLen; ix++ {
		startPc, err := utils.ReadInt16(reader)
		if nil != err {
			return nil, err
		}

		endPc, err := utils.ReadInt16(reader)
		if nil != err {
			return nil, err
		}

		handlerPc, err := utils.ReadInt16(reader)
		if nil != err {
			return nil, err
		}

		catchType, err := utils.ReadInt16(reader)
		if nil != err {
			return nil, err
		}

		expTable := &ExceptionTable{
			StartPc:   startPc,
			EndPc:     endPc,
			HandlerPc: handlerPc,
			CatchType: catchType,
		}

		attr.ExceptionTable = append(attr.ExceptionTable, expTable)

	}

	// 属性长度
	attrTot, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, err
	}

	// 读取属性表
	attr.Attrs = make([]interface{}, 0, attrTot)
	for ix := 0; ix < int(attrTot); ix++ {
		at, err := c.ReadAttr(reader)
		if nil != err {
			return nil, fmt.Errorf("failed to read attr, %w", err)
		}

		attr.Attrs = append(attr.Attrs, at)
	}

	return attr, nil
}

type ConstantValueAttr struct {
	AttrLength uint32

	ConstantValueIndex uint16
}

func (c *ConstantValueAttr) String() string {
	return "ConstantValue"
}

func ReadConstantValueAttr(reader io.Reader) (*ConstantValueAttr, error) {
	length, err := utils.ReadInt32(reader)
	if nil != err {
		return nil, err
	}

	valIndex, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, err
	}

	return &ConstantValueAttr{
		AttrLength:         length,
		ConstantValueIndex: valIndex,
	}, nil
}

type ExceptionAttr struct {
	AttrNameIndex uint16
	AttrLength uint32

	NumberOfExceptions uint16
	ExceptionIndexTable []uint16
}

func (e *ExceptionAttr) String() string {
	return "Exception"
}

func ReadExceptionAttr(reader io.Reader) (*ExceptionAttr, error) {
	panic("implement me")
}

type LineNumberAttr struct {
	// AttrNameIndex uint16
	AttrLength uint32

	LineNumberTableLength uint16
	LineNumberTable []*LineNumberInfo
}

func (l *LineNumberAttr) String() string {
	return "LineNumber"
}

type LineNumberInfo struct {
	StartPc uint16
	LineNumber uint16
}

func ReadLineNumberAttr(reader io.Reader) (*LineNumberAttr, error) {
	//nameIndex, err := utils.ReadInt16(reader)
	//if nil != err {
	//	return nil, fmt.Errorf("failed to load attr_name_index, %w", err)
	//}

	length, err := utils.ReadInt32(reader)
	if nil != err {
		return nil, fmt.Errorf("failed to load attr_length, %w", err)
	}

	tableLen, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, fmt.Errorf("failed to load line_number_table_length, %w", err)
	}

	tableLength := int(tableLen)
	tables := make([]*LineNumberInfo, 0, tableLength)
	for ix := 0; ix < tableLength; ix++ {
		startPc, err := utils.ReadInt16(reader)
		if nil != err {
			return nil, fmt.Errorf("failed to load start_pc, %w", err)
		}

		lineNumber, err := utils.ReadInt16(reader)
		if nil != err {
			return nil, fmt.Errorf("failed to load line_number, %w", err)
		}


		table := &LineNumberInfo{
			StartPc:    startPc,
			LineNumber: lineNumber,
		}
		tables = append(tables, table)
	}

	return &LineNumberAttr{
		//AttrNameIndex:         nameIndex,
		AttrLength:            length,
		LineNumberTableLength: tableLen,
		LineNumberTable:       tables,
	}, nil
}

type LocalVariableAttr struct {
	AttrNameIndex uint16
	AttrLength uint32

	LocalVariableLength uint16
	LocalVariableInfo []*LocalVariableTable
}

func (l *LocalVariableAttr) String() string {
	return "LovalVariable"
}

type LocalVariableTable struct {
	StartPc uint16
	Length uint16
	NameIndex uint16
	DescriptorIndex uint16
	Index uint16
}

func ReadLocalVariableAttr(reader io.Reader) (*LocalVariableAttr, error) {
	nameIndex, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, fmt.Errorf("failed to load attr_name_index, %w", err)
	}

	length, err := utils.ReadInt32(reader)
	if nil != err {
		return nil, fmt.Errorf("failed to load attr_length, %w", err)
	}

	varLen, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, fmt.Errorf("failed to load local_variable_length, %w", err)
	}

	varLength := int(varLen)
	tables := make([]*LocalVariableTable, 0, varLength)
	for ix := 0; ix < varLength; ix++ {
		startPc, err := utils.ReadInt16(reader)
		if nil != err {
			return nil, fmt.Errorf("failed to load start_pc, %w", err)
		}

		length, err := utils.ReadInt16(reader)
		if nil != err {
			return nil, fmt.Errorf("failed to load length, %w", err)
		}

		tableNameIndex, err := utils.ReadInt16(reader)
		if nil != err {
			return nil, fmt.Errorf("failed to load name_index, %w", err)
		}

		despIndex, err := utils.ReadInt16(reader)
		if nil != err {
			return nil, fmt.Errorf("failed to load descriptor_index, %w", err)
		}

		idx, err := utils.ReadInt16(reader)
		if nil != err {
			return nil, fmt.Errorf("failed to load index, %w", err)
		}

		table := &LocalVariableTable{
			StartPc:         startPc,
			Length:          length,
			NameIndex:       tableNameIndex,
			DescriptorIndex: despIndex,
			Index:           idx,
		}

		tables = append(tables, table)
	}

	return &LocalVariableAttr{
		AttrNameIndex:       nameIndex,
		AttrLength:          length,
		LocalVariableLength: varLen,
		LocalVariableInfo:   tables,
	}, nil
}



type InnerClassAttr struct {
	Length uint32
	NumberOfClasses uint16
	InnerClasses []*InnerClassInfo
}

type InnerClassInfo struct {
	InnerClassInfoIndex uint16
	OuterClassInfoIndex uint16
	InnerNameIndex uint16
	InnerClassAccessFlags uint16
}

func ReadInnerClassAttr(reader io.Reader) (*InnerClassAttr, error) {
	length, err := utils.ReadInt32(reader)
	if nil != err {
		return nil, fmt.Errorf("failed to load length: %w", err)
	}

	numberOfClasses, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, fmt.Errorf("failed to load number_of_classes: %w", err)
	}

	innerInfos := make([]*InnerClassInfo, 0, numberOfClasses)
	for ix := 0; ix < int(numberOfClasses); ix++ {
		innerIndex, err := utils.ReadInt16(reader)
		if nil != err {
			return nil, fmt.Errorf("failed to load inner_class_info_index: %w", err)
		}

		outerIndex, err := utils.ReadInt16(reader)
		if nil != err {
			return nil, fmt.Errorf("failed to load outer_class_info_index: %w", err)
		}

		nameIndex, err := utils.ReadInt16(reader)
		if nil != err {
			return nil, fmt.Errorf("failed to load inner_name_index: %w", err)
		}

		accessFlags, err := utils.ReadInt16(reader)
		if nil != err {
			return nil, fmt.Errorf("failed to load inner_class_access_flags: %w", err)
		}

		info := &InnerClassInfo{
			InnerClassInfoIndex:   innerIndex,
			OuterClassInfoIndex:   outerIndex,
			InnerNameIndex:        nameIndex,
			InnerClassAccessFlags: accessFlags,
		}

		innerInfos = append(innerInfos, info)
	}

	return &InnerClassAttr{
		Length:          length,
		NumberOfClasses: numberOfClasses,
		InnerClasses:    innerInfos,
	}, nil
}

type SourceFileAttr struct {
	Length uint32
	SourceFileIndex uint16
}

func (s *SourceFileAttr) String() string {
	return "SourceFile"
}


func ReadSourceFileAttr(reader io.Reader) (*SourceFileAttr, error) {
	length, err := utils.ReadInt32(reader)
	if nil != err {
		return nil, fmt.Errorf("failed to load length: %w", err)
	}

	idx, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, fmt.Errorf("failed to load source_file_index: %w", err)
	}

	return &SourceFileAttr{
		Length:          length,
		SourceFileIndex: idx,
	}, nil
}
