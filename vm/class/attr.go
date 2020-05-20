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
			return nil, fmt.Errorf("failed to read line_number_table attr: %w", err)
		}

		return lineAttr, nil

	} else if "SourceFile" == attrName {
		srcAttr, err := ReadSourceFileAttr(reader)
		if nil != err {
			return nil, fmt.Errorf("failed to read source_file attr: %w", err)
		}

		return srcAttr, nil

	} else if "StackMapTable" == attrName {
		// 跳过此属性
		err := c.skipAttr(reader)
		if nil != err {
			return nil, fmt.Errorf("failed to skip StackMapTable attr: %w", err)
		}

		return struct{}{}, nil
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
