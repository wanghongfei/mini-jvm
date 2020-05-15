package class

import (
	"github.com/wanghongfei/mini-jvm/utils"
	"io"
)

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

