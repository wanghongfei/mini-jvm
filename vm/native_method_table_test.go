package vm

import (
	"github.com/wanghongfei/mini-jvm/vm/class"
	"testing"
)

func TestNativeMethodInfo_ParseArgCount(t *testing.T) {
	info := &NativeMethodInfo{
		Name:       "",
		Descriptor: "(II)V",
		EntryFunc:  nil,
	}
	if 2 != class.ParseArgCount(info.Descriptor) {
		t.FailNow()
	}

	info = &NativeMethodInfo{
		Name:       "",
		Descriptor: "(I)V",
		EntryFunc:  nil,
	}
	if 1 != class.ParseArgCount(info.Descriptor) {
		t.FailNow()
	}

	info = &NativeMethodInfo{
		Name:       "",
		Descriptor: "()V",
		EntryFunc:  nil,
	}
	if 0 != class.ParseArgCount(info.Descriptor) {
		t.FailNow()
	}

	info = &NativeMethodInfo{
		Name:       "",
		Descriptor: "(Ljava/lang/Runnable;)V",
		EntryFunc:  nil,
	}
	if 1 != class.ParseArgCount(info.Descriptor) {
		t.FailNow()
	}

	info = &NativeMethodInfo{
		Name:       "",
		Descriptor: "(Ljava/lang/Runnable;II)V",
		EntryFunc:  nil,
	}
	if 3 != class.ParseArgCount(info.Descriptor) {
		t.FailNow()
	}

	info = &NativeMethodInfo{
		Name:       "",
		Descriptor: "(CCLjava/lang/Runnable;II)V",
		EntryFunc:  nil,
	}
	if 5 != class.ParseArgCount(info.Descriptor) {
		t.FailNow()
	}
}
