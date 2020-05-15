package class

import (
	"github.com/wanghongfei/mini-jvm/utils"
	"io"
)

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
