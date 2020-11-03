package vm

import "github.com/wanghongfei/mini-jvm/vm/class"

// 遇到athrow指令, 当前方法的异常表中匹配不到异常时返回此错误
type ExceptionThrownError struct {
	ExceptionRef *class.Reference
}

func (e ExceptionThrownError) Error() string {
	return "throw exception: " + e.ExceptionRef.Object.DefFile.FullClassName
}

func NewExceptionThrownError(ref *class.Reference) error {
	return &ExceptionThrownError{ExceptionRef: ref}
}

