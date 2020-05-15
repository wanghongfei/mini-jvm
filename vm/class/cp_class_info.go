package class

import (
	"github.com/wanghongfei/mini-jvm/utils"
	"io"
)

type ClassInfoConstInfo struct {
	Tag uint8
	FullClassNameIndex uint16
}

func ReadClassInfoConst(reader io.Reader, tag uint8) (*ClassInfoConstInfo, error) {
	cpInfo := new(ClassInfoConstInfo)
	cpInfo.Tag = tag

	classIndex, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, err
	}
	cpInfo.FullClassNameIndex = classIndex

	return cpInfo, nil
}

