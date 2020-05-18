package class

import (
	"github.com/wanghongfei/mini-jvm/utils"
	"io"
)

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

