package class

import (
	"github.com/wanghongfei/mini-jvm/utils"
	"io"
)

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

