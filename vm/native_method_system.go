package vm

import "github.com/wanghongfei/mini-jvm/vm/class"

//public static native void arraycopy(Object src,  int  srcPos,
//                                    Object dest, int destPos,
//									  int length);
func SystemArrayCopy(args ...interface{}) interface{} {
	rawSrc := args[2]
	rawSrcPos := args[3]
	rawDest := args[4]
	rawDestPos := args[5]
	rawLength := args[6]

	srcArr := rawSrc.(*class.Reference).Array.Data
	srcPos := rawSrcPos.(int)
	destArr := rawDest.(*class.Reference).Array.Data
	destPos := rawDestPos.(int)
	length := rawLength.(int)

	for ix := 0; ix < length; ix++ {
		destArr[destPos + ix] = srcArr[srcPos + ix]
	}

	return nil
}
