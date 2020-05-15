package class

import (
	"fmt"
	"github.com/wanghongfei/mini-jvm/utils"
	"io"
)

type SourceFileAttr struct {
	Length uint32
	SourceFileIndex uint16
}

func (s *SourceFileAttr) String() string {
	return "SourceFile"
}

func ReadSourceFileAttr(reader io.Reader) (*SourceFileAttr, error) {
	length, err := utils.ReadInt32(reader)
	if nil != err {
		return nil, fmt.Errorf("failed to load length: %w", err)
	}

	idx, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, fmt.Errorf("failed to load source_file_index: %w", err)
	}

	return &SourceFileAttr{
		Length:          length,
		SourceFileIndex: idx,
	}, nil
}
