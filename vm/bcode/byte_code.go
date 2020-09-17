package bcode

import "encoding/hex"

const (
	Nop byte = 0x00
	Aconstnull = 0x01

	Iconst0 = 0x03
	Iconst1 = 0x04
	Iconst2 = 0x05
	Iconst3 = 0x06
	Iconst4 = 0x07
	Iconst5 = 0x08

	Ldc = 0x12

	Iaload = 0x2e

	Aaload = 0x32
	Caload = 0x34

	Istore0 = 0x3b
	Istore1 = 0x3c
	Istore2 = 0x3d
	Istore3 = 0x3e

	Bipush = 0x10
	Sipush = 0x11

	Iload = 0x15
	Iload0 = 0x1a
	Iload1 = 0x1b
	Iload2 = 0x1c
	Iload3 = 0x1d

	Aload = 0x19
	Aload0 = 0x2a
	Aload1 = 0x2b
	Aload2 = 0x2c
	Aload3 = 0x2d

	Getstatic = 0xb2
	Putstatic = 0xb3

	Athrow = 0xbf

	Monitorenter = 0xc2
	Monitorexit = 0xc3

	Istore = 0x36
	Lstore1 = 0x40

	Astore = 0x3a
	Astore0 = 0x4b
	Astore1 = 0x4c
	Astore2 = 0x4d
	Astore3 = 0x4e
	Iastore = 0x4f

	Aastore = 0x53
	Castore = 0x55
	Pop = 0x57

	Dup = 0x59

	Iadd = 0x60
	Isub = 0x64

	Ishl = 0x78

	Iinc = 0x84

	Ifeq = 0x99
	Ifne = 0x9a
	Iflt = 0x9b
	Ifge = 0x9c
	Ifgt = 0x9d
	Ifle = 0x9e

	Ificmpeq = 0x9f
	Ificmpne = 0xa0
	Ificmplt = 0xa1
	Ificmpge = 0xa2
	Ificmpgt = 0xa3
	Ificmple = 0xa4
	Ifacmpeq = 0xa5
	Ifacmpne = 0xa6
	Goto = 0xa7

	Areturn = 0xb0
	Return = 0xb1

	GetField = 0xb4
	Putfield = 0xb5

	Newarray = 0xbc
	Anewarray = 0xbd

	Invokevirtual = 0xb6
	Invokespecial = 0xb7
	Invokestatic = 0xb8
	Invokeinterface = 0xb9

	New = 0xbb

	Arraylength = 0xbe

	Ireturn = 0xac

	Wide = 0xc4
	Ifnonnull = 0xc7
)

func ToName(code byte) string {
	switch code {
	case Aconstnull:
		return "aconstnull"

	case Iconst0:
		return "iconst_0"
	case Iconst1:
		return "iconst_1"
	case Iconst2:
		return "iconst_2"
	case Iconst3:
		return "iconst_3"
	case Iconst4:
		return "iconst_4"
	case Iconst5:
		return "iconst_5"

	case Ldc:
		return "ldc"

	case Iaload:
		return "iaload"
	case Aaload:
		return "aaload"
	case Caload:
		return "caload"

	case Istore0:
		return "istore_0"
	case Istore1:
		return "istore_1"
	case Istore2:
		return "istore_2"
	case Istore3:
		return "istore_3"

	case Bipush:
		return "bipush"
	case Sipush:
		return "sipush"

	case Iload:
		return "iload"
	case Iload0:
		return "iload_0"
	case Iload1:
		return "iload_1"
	case Iload2:
		return "iload_2"
	case Iload3:
		return "iload_3"

	case Aload:
		return "aload"
	case Aload0:
		return "aload_0"
	case Aload1:
		return "aload_1"
	case Aload2:
		return "aload_2"
	case Aload3:
		return "aload_3"

	case Getstatic:
		return "getstatic"
	case Putstatic:
		return "putstatic"

	case Athrow:
		return "athrow"

	case Monitorenter:
		return "monitorenter"
	case Monitorexit:
		return "monitorexit"

	case Istore:
		return "istore"

	case Lstore1:
		return "lstore_1"

	case Astore:
		return "astore"
	case Astore0:
		return "astore_0"
	case Astore1:
		return "astore_1"
	case Astore2:
		return "astore_2"
	case Astore3:
		return "astore_3"

	case Iastore:
		return "iastore"
	case Aastore:
		return "aastore"
	case Castore:
		return "castore"

	case Pop:
		return "pop"
	case Dup:
		return "dup"

	case Iadd:
		return "iadd"
	case Isub:
		return "isub"
	case Ishl:
		return "ishl"
	case Iinc:
		return "iinc"

	case Ifeq:
		return "ifeq"
	case Ifne:
		return "ifne"
	case Iflt:
		return "iflt"
	case Ifge:
		return "ifge"
	case Ifgt:
		return "ifgt"
	case Ifle:
		return "ifle"
	case Ificmpeq:
		return "ificmpeq"
	case Ificmpne:
		return "ificmpne"
	case Ificmplt:
		return "ificmplt"
	case Ificmpgt:
		return "ificmpgt"
	case Ificmpge:
		return "ificmpge"
	case Ificmple:
		return "ificmple"
	case Ifacmpeq:
		return "ifacmpeq"
	case Ifacmpne:
		return "ifacmpne"

	case Goto:
		return "goto"

	case Areturn:
		return "areturn"
	case Return:
		return "return"

	case GetField:
		return "getfield"
	case Putfield:
		return "putfield"

	case Newarray:
		return "newarray"
	case Anewarray:
		return "anewarray"

	case Invokevirtual:
		return "invokevirtual"
	case Invokespecial:
		return "invokespecial"
	case Invokestatic:
		return "invokestatic"
	case Invokeinterface:
		return "invokeinterface"

	case New:
		return "new"
	case Arraylength:
		return "arraylength"

	case Ireturn:
		return "ireturn"

	case Wide:
		return "wide"

	case Ifnonnull:
		return "ifnonnull"

	default:
		return "unknown: " + hex.EncodeToString([]byte{code})
	}
}
