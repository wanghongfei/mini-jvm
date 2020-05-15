package class

import (
	"fmt"
	"github.com/wanghongfei/mini-jvm/utils"
	"io"
)

type MethodInfo struct {
	AccessFlags uint16
	NameIndex uint16
	DescriptorIndex uint16

	AttrCount uint16
	Attrs []interface{}

	// 所在的class定义文件
	DefFile *DefFile
}

func (f *MethodInfo) String() string {
	return f.DefFile.ConstPool[f.NameIndex].(*Utf8InfoConst).String()
}

func (c *DefFile) ReadMethodInfo(reader io.Reader) (*MethodInfo, error) {
	flags, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, fmt.Errorf("failed to read access_flags, %w", err)
	}

	nameIndex, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, fmt.Errorf("failed to read name_index, %w", err)
	}

	descriptorIndex, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, fmt.Errorf("failed to read descriptor_index, %w", err)
	}

	attrCount, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, fmt.Errorf("failed to read attr_count, %w", err)
	}

	// 读取属性表
	attrs := make([]interface{}, 0, attrCount)
	for ix := 0; ix < int(attrCount); ix++ {
		attr, err := c.ReadAttr(reader)
		if nil != err {
			return nil, fmt.Errorf("failed to read attr, %w", err)
		}

		attrs = append(attrs, attr)
	}

	return &MethodInfo{
		DefFile:         c,
		AccessFlags:     flags,
		NameIndex:       nameIndex,
		DescriptorIndex: descriptorIndex,
		AttrCount:       attrCount,
		Attrs:           attrs,
	}, nil
}
