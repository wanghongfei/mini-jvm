package class

import (
	"fmt"
	"github.com/wanghongfei/mini-jvm/utils"
	"io"
)

// code属性
type CodeAttr struct {
	AttrLength uint32

	MaxStack uint16
	MaxLocals uint16

	// 字节码长度
	CodeLength uint32
	Code []byte

	// 异常表
	ExceptionTableLength uint16
	ExceptionTable []*ExceptionTable

	AttrCount uint16
	Attrs []interface{}
}

func (c *CodeAttr) String() string {
	return "Code"
}

type ExceptionTable struct {
	StartPc uint16
	EndPc uint16
	HandlerPc uint16
	CatchType uint16
}

func (c *DefFile) ReadCodeAttr(reader io.Reader) (*CodeAttr, error) {
	attr := new(CodeAttr)

	length, err := utils.ReadInt32(reader)
	if nil != err {
		return nil, err
	}
	attr.AttrLength = length

	maxStack, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, err
	}
	attr.MaxStack = maxStack

	maxLocals, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, err
	}
	attr.MaxLocals = maxLocals

	// 字节码长度
	codeLen, err := utils.ReadInt32(reader)
	if nil != err {
		return nil, err
	}
	attr.CodeLength = codeLen

	// 字节码
	codeBuf := make([]byte, codeLen)
	_, err = reader.Read(codeBuf)
	if nil != err {
		return nil, err
	}
	attr.Code = codeBuf

	// 异常表长度
	expTableLen, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, err
	}
	attr.ExceptionTableLength = expTableLen

	// 异常表
	expLen := int(expTableLen)
	attr.ExceptionTable = make([]*ExceptionTable, 0, expLen)
	for ix := 0; ix < expLen; ix++ {
		startPc, err := utils.ReadInt16(reader)
		if nil != err {
			return nil, err
		}

		endPc, err := utils.ReadInt16(reader)
		if nil != err {
			return nil, err
		}

		handlerPc, err := utils.ReadInt16(reader)
		if nil != err {
			return nil, err
		}

		catchType, err := utils.ReadInt16(reader)
		if nil != err {
			return nil, err
		}

		expTable := &ExceptionTable{
			StartPc:   startPc,
			EndPc:     endPc,
			HandlerPc: handlerPc,
			CatchType: catchType,
		}

		attr.ExceptionTable = append(attr.ExceptionTable, expTable)

	}

	// 属性长度
	attrTot, err := utils.ReadInt16(reader)
	if nil != err {
		return nil, err
	}

	// 读取属性表
	attr.Attrs = make([]interface{}, 0, attrTot)
	for ix := 0; ix < int(attrTot); ix++ {
		at, err := c.ReadAttr(reader)
		if nil != err {
			return nil, fmt.Errorf("failed to read attr, %w", err)
		}

		attr.Attrs = append(attr.Attrs, at)
	}

	return attr, nil
}
