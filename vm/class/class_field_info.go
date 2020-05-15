package class

import (
	"fmt"
	"github.com/wanghongfei/mini-jvm/utils"
	"io"
)

type FieldInfo struct {
	AccessFlags uint16
	NameIndex uint16
	DescriptorIndex uint16

	AttrCount uint16
	Attrs []interface{}
}

func (c *DefFile) ReadFieldInfo(reader io.Reader) (*FieldInfo, error) {
	info := new(FieldInfo)

	accessFlags, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, err
	}
	info.AccessFlags = accessFlags

	nameIndex, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, err
	}
	info.NameIndex = nameIndex

	descIndex, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, err
	}
	info.DescriptorIndex = descIndex

	attrCount, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, err
	}
	info.AttrCount = attrCount

	info.Attrs = make([]interface{}, 0, attrCount)
	for ix := 0; ix < int(attrCount); ix++ {
		attr, err := c.ReadAttr(reader)
		if nil != err {
			return nil, fmt.Errorf("failed to read attr, %w", err)
		}

		info.Attrs = append(info.Attrs, attr)
	}

	return info, nil
}
