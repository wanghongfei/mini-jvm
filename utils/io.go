package utils

import (
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
	"os"
)

func ReadAllFromFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if nil != err {
		return nil, err
	}
	defer f.Close()

	return ioutil.ReadAll(f)
}

func ReadInt32(bufReader io.Reader) (uint32, error) {
	numBuf := make([]byte, 4, 4)
	_, err := bufReader.Read(numBuf)
	if nil != err {
		return 0, err
	}

	var num uint32
	err = binary.Read(bytes.NewBuffer(numBuf), binary.BigEndian, &num)
	if nil != err {
		return 0, err
	}

	return num, nil
}


func ReadInt16(bufReader io.Reader) (uint16, error) {
	numBuf := make([]byte, 2, 2)
	_, err := bufReader.Read(numBuf)
	if nil != err {
		return 0, err
	}

	var num uint16
	err = binary.Read(bytes.NewBuffer(numBuf), binary.BigEndian, &num)
	if nil != err {
		return 0, err
	}

	return num, nil
}

func ReadInt8(bufReader io.Reader) (uint8, error) {
	numBuf := make([]byte, 1, 1)
	_, err := bufReader.Read(numBuf)
	if nil != err {
		return 0, err
	}

	var num uint8
	err = binary.Read(bytes.NewBuffer(numBuf), binary.BigEndian, &num)
	if nil != err {
		return 0, err
	}

	return num, nil
}
