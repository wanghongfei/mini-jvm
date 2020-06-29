package vm

import "testing"

func TestNativeMethodInfo_ParseArgCount(t *testing.T) {
	info := &NativeMethodInfo{
		Name:       "",
		Descriptor: "(II)V",
		EntryFunc:  nil,
	}
	if 2 != info.ParseArgCount() {
		t.FailNow()
	}

	info = &NativeMethodInfo{
		Name:       "",
		Descriptor: "(I)V",
		EntryFunc:  nil,
	}
	if 1 != info.ParseArgCount() {
		t.FailNow()
	}

	info = &NativeMethodInfo{
		Name:       "",
		Descriptor: "()V",
		EntryFunc:  nil,
	}
	if 0 != info.ParseArgCount() {
		t.FailNow()
	}

	info = &NativeMethodInfo{
		Name:       "",
		Descriptor: "(Ljava/lang/Runnable;)V",
		EntryFunc:  nil,
	}
	if 1 != info.ParseArgCount() {
		t.FailNow()
	}

	info = &NativeMethodInfo{
		Name:       "",
		Descriptor: "(Ljava/lang/Runnable;II)V",
		EntryFunc:  nil,
	}
	if 3 != info.ParseArgCount() {
		t.FailNow()
	}

	info = &NativeMethodInfo{
		Name:       "",
		Descriptor: "(CCLjava/lang/Runnable;II)V",
		EntryFunc:  nil,
	}
	if 5 != info.ParseArgCount() {
		t.FailNow()
	}
}
