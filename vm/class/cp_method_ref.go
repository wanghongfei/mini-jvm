package class

import (
	"github.com/wanghongfei/mini-jvm/utils"
	"io"
)

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
