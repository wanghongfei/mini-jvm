package class

import (
	"github.com/wanghongfei/mini-jvm/utils"
	"io"
)

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

