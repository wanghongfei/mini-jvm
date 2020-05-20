package vm

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/wanghongfei/mini-jvm/vm/class"
)

// 解释执行引擎
type InterpretedExecutionEngine struct {
	miniJvm *MiniJvm

	methodStack *MethodStack
}


func (i *InterpretedExecutionEngine) Execute(def *class.DefFile, methodName string) error {
	// 查找方法
	method, err := i.findMethod(def, methodName)
	if nil != err {
		return fmt.Errorf("failed to find method: %w", err)
	}

	// 提取code属性
	codeAttr, err := i.findCodeAttr(method)
	if nil != err {
		return fmt.Errorf("failed to extract code attr: %w", err)
	}

	// 创建栈帧
	frame := newMethodStackFrame(int(codeAttr.MaxStack), int(codeAttr.MaxLocals))
	// 压栈
	i.methodStack.Push(frame)

	for {
		// 取出pc指向的字节码
		byteCode := codeAttr.Code[frame.pc]

		exitLoop := false

		// 执行
		switch byteCode {
		case iconst0:
			// 将x压栈
			frame.opStack.Push(0)
		case iconst1:
			frame.opStack.Push(1)
		case iconst2:
			frame.opStack.Push(2)


		case istore1:
			// 将栈顶int型数值存入第二个本地变量
			top, _ := frame.opStack.Pop()
			frame.localVariablesTable[1] = top
		case istore2:
			// 将栈顶int型数值存入第3个本地变量
			top, _ := frame.opStack.Pop()
			frame.localVariablesTable[2] = top
		case istore3:
			// 将栈顶int型数值存入第4个本地变量
			top, _ := frame.opStack.Pop()
			frame.localVariablesTable[3] = top

		case iload1:
			// 将第1个slot中的值压栈
			frame.opStack.Push(frame.localVariablesTable[1])
		case iload2:
			frame.opStack.Push(frame.localVariablesTable[2])

		case iadd:
			// 取出栈顶2元素，相加，入栈
			op1, _ := frame.opStack.Pop()
			op2, _ := frame.opStack.Pop()
			sum := op1 + op2
			frame.opStack.Push(sum)

		case bipush:
			// 将单字节的常量值(-128~127)推送至栈顶
			num := codeAttr.Code[frame.pc + 1]
			frame.opStack.Push(uint32(num))
			frame.pc++

		case ificmpgt:
			// 比较栈顶两int型数值大小, 当结果大于0时跳转

			// 待比较的数
			x, _ := frame.opStack.Pop()
			y, _ := frame.opStack.Pop()

			// 跳转的偏移量
			twoByteNum := codeAttr.Code[frame.pc + 1 : frame.pc + 1 + 2]
			var offset int16
			err := binary.Read(bytes.NewBuffer(twoByteNum), binary.BigEndian, &offset)
			if nil != err {
				return fmt.Errorf("failed to read offset for if_icmpgt: %w", err)
			}

			if int(y) - int(x) > 0 {
				frame.pc = frame.pc + int(offset) - 1

			} else {
				frame.pc += 2
			}


		case iinc:
			// 将第op1个slot的变量增加op2
			op1 := codeAttr.Code[frame.pc + 1]
			op2 := codeAttr.Code[frame.pc + 2]
			frame.pc += 2

			frame.localVariablesTable[op1] = frame.localVariablesTable[op1] + uint32(op2)

		case bgoto:
			// 跳转
			twoByteNum := codeAttr.Code[frame.pc + 1 : frame.pc + 1 + 2]
			var offset int16
			err := binary.Read(bytes.NewBuffer(twoByteNum), binary.BigEndian, &offset)
			if nil != err {
				return fmt.Errorf("failed to read pc offset for 'goto': %w", err)
			}

			frame.pc = frame.pc + int(offset) - 1

		case emptyreturn:
			// 返回
			exitLoop = true

		default:
			return fmt.Errorf("unsupported byte code %s", hex.EncodeToString([]byte{byteCode}))
		}

		if exitLoop {
			break
		}

		// 移动程序计数器
		frame.pc++
	}


	i.methodStack.Pop()

	return nil
}

func (i *InterpretedExecutionEngine) findCodeAttr(method *class.MethodInfo) (*class.CodeAttr, error) {
	for _, attrGeneric := range method.Attrs {
		attr, ok := attrGeneric.(*class.CodeAttr)
		if ok {
			return attr, nil
		}
	}

	return nil, errors.New("no node attr in method")
}

func (i *InterpretedExecutionEngine) findMethod(def *class.DefFile, methodName string) (*class.MethodInfo, error) {
	for _, method := range def.Methods {
		name := def.ConstPool[method.NameIndex].(*class.Utf8InfoConst).String()
		if name == methodName {
			return method, nil
		}
	}

	return nil, fmt.Errorf("method '%s' not found", methodName)
}

func NewInterpretedExecutionEngine(vm *MiniJvm) *InterpretedExecutionEngine {
	return &InterpretedExecutionEngine{
		miniJvm:     vm,
		methodStack: NewMethodStack(1024),
	}
}

