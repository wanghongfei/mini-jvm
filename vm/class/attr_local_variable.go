package class

import (
	"fmt"
	"github.com/wanghongfei/mini-jvm/utils"
	"io"
)

type LocalVariableAttr struct {
	AttrNameIndex uint16
	AttrLength uint32

	LocalVariableLength uint16
	LocalVariableInfo []*LocalVariableTable
}

func (l *LocalVariableAttr) String() string {
	return "LovalVariable"
}

type LocalVariableTable struct {
	StartPc uint16
	Length uint16
	NameIndex uint16
	DescriptorIndex uint16
	Index uint16
}

func ReadLocalVariableAttr(reader io.Reader) (*LocalVariableAttr, error) {
	nameIndex, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, fmt.Errorf("failed to load attr_name_index, %w", err)
	}

	length, err := utils.ReadInt32(reader)
	if nil != err {
		return nil, fmt.Errorf("failed to load attr_length, %w", err)
	}

	varLen, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, fmt.Errorf("failed to load local_variable_length, %w", err)
	}

	varLength := int(varLen)
	tables := make([]*LocalVariableTable, 0, varLength)
	for ix := 0; ix < varLength; ix++ {
		startPc, err := utils.ReadInt16(reader)
		if nil != err {
			return nil, fmt.Errorf("failed to load start_pc, %w", err)
		}

		length, err := utils.ReadInt16(reader)
		if nil != err {
			return nil, fmt.Errorf("failed to load length, %w", err)
		}

		tableNameIndex, err := utils.ReadInt16(reader)
		if nil != err {
			return nil, fmt.Errorf("failed to load name_index, %w", err)
		}

		despIndex, err := utils.ReadInt16(reader)
		if nil != err {
			return nil, fmt.Errorf("failed to load descriptor_index, %w", err)
		}

		idx, err := utils.ReadInt16(reader)
		if nil != err {
			return nil, fmt.Errorf("failed to load index, %w", err)
		}

		table := &LocalVariableTable{
			StartPc:         startPc,
			Length:          length,
			NameIndex:       tableNameIndex,
			DescriptorIndex: despIndex,
			Index:           idx,
		}

		tables = append(tables, table)
	}

	return &LocalVariableAttr{
		AttrNameIndex:       nameIndex,
		AttrLength:          length,
		LocalVariableLength: varLen,
		LocalVariableInfo:   tables,
	}, nil
}

