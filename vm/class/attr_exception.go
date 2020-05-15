package class

import "io"

type ExceptionAttr struct {
	AttrNameIndex uint16
	AttrLength uint32

	NumberOfExceptions uint16
	ExceptionIndexTable []uint16
}


func ReadExceptionAttr(reader io.Reader) (*ExceptionAttr, error) {
	panic("implement me")
}