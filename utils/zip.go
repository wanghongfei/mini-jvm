package utils

import (
	"archive/zip"
	"fmt"
	"io"
)

// zip内部文件访问函数
type ZipFileVisitor func(reader io.Reader) (bool, error)
// 判断当前zip文件是否为目标文件的函数
type ZipFilePredicate func(f *zip.File) bool

func VisitZip(zipFile string, predicate ZipFilePredicate, visitor ZipFileVisitor) error {
	zr, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zr.Close()

	for _, file := range zr.File {
		// 判断用户是否对此文件感兴趣
		if !predicate(file) {
			continue
		}

		// 打开文件
		innerFile, err := file.Open()
		if nil != err {
			return fmt.Errorf("failed to open inner file '%s': %w", file.Name, err)
		}

		// 调用业务逻辑
		stop, err := visitor(innerFile)
		if nil != err {
			innerFile.Close()
			return fmt.Errorf("biz function returned error: %w", err)
		}

		// 停止遍历后续文件
		if stop {
			innerFile.Close()
			return nil
		}

		innerFile.Close()
	}

	return nil
}
