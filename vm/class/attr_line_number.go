package class

import (
	"fmt"
	"github.com/wanghongfei/mini-jvm/utils"
	"io"
)

type LineNumberAttr struct {
	// AttrNameIndex uint16
	AttrLength uint32

	LineNumberTableLength uint16
	LineNumberTable []*LineNumberInfo
}

type LineNumberInfo struct {
	StartPc uint16
	LineNumber uint16
}

func ReadLineNumberAttr(reader io.Reader) (*LineNumberAttr, error) {
	//nameIndex, err := utils.ReadInt16(reader)
	//if nil != err {
	//	return nil, fmt.Errorf("failed to load attr_name_index, %w", err)
	//}

	length, err := utils.ReadInt32(reader)
	if nil != err {
		return nil, fmt.Errorf("failed to load attr_length, %w", err)
	}

	tableLen, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, fmt.Errorf("failed to load line_number_table_length, %w", err)
	}

	tableLength := int(tableLen)
	tables := make([]*LineNumberInfo, 0, tableLength)
	for ix := 0; ix < tableLength; ix++ {
		startPc, err := utils.ReadInt16(reader)
		if nil != err {
			return nil, fmt.Errorf("failed to load start_pc, %w", err)
		}

		lineNumber, err := utils.ReadInt16(reader)
		if nil != err {
			return nil, fmt.Errorf("failed to load line_number, %w", err)
		}


		table := &LineNumberInfo{
			StartPc:    startPc,
			LineNumber: lineNumber,
		}
		tables = append(tables, table)
	}

	return &LineNumberAttr{
		//AttrNameIndex:         nameIndex,
		AttrLength:            length,
		LineNumberTableLength: tableLen,
		LineNumberTable:       tables,
	}, nil
}