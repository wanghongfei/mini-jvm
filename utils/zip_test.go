package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"testing"
)

func TestVisitZip(t *testing.T) {
	VisitZip("/Library/Java/JavaVirtualMachines/jdk1.8.0_181.jdk/Contents/Home/jre/lib/rt.jar", func(f *zip.File) bool {
		return f.Name == "java/lang/String.class"

	}, func(reader io.Reader) (bool, error) {
		fmt.Println("Got You!")
		return true, nil
	})
}

