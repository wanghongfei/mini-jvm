package class

import (
	"github.com/wanghongfei/mini-jvm/utils"
	"io"
)

type Utf8InfoConst struct {
	Tag uint8
	Length uint16
	Bytes []byte
}

func ReadUtf8InfoConst(reader io.Reader, tag uint8) (*Utf8InfoConst, error) {
	cpInfo := new(Utf8InfoConst)
	cpInfo.Tag = tag

	length, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, err
	}
	cpInfo.Length = length

	utf8Buf := make([]byte, length)
	_, err = reader.Read(utf8Buf)
	if nil != err {
		return nil, err
	}
	cpInfo.Bytes = utf8Buf

	return cpInfo, nil
}

func (o *Utf8InfoConst) String() string {
	return string(o.Bytes)
}

