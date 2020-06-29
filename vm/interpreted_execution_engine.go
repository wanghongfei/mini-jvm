package vm

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/wanghongfei/mini-jvm/vm/accflag"
	"github.com/wanghongfei/mini-jvm/vm/bcode"
	"github.com/wanghongfei/mini-jvm/vm/class"
)

// 解释执行引擎
type InterpretedExecutionEngine struct {
	miniJvm *MiniJvm

	// methodStack *MethodStack
}

func (i *InterpretedExecutionEngine) Execute(def *class.DefFile, methodName string) error {
	return i.ExecuteWithFrame(def, methodName, "([Ljava/lang/String;)V", nil)
}

func (i *InterpretedExecutionEngine) ExecuteWithDescriptor(def *class.DefFile, methodName, descriptor string) error {
	return i.ExecuteWithFrame(def, methodName, descriptor, nil)
}

func (i *InterpretedExecutionEngine) ExecuteWithFrame(def *class.DefFile, methodName string, methodDescriptor string, lastFrame *MethodStackFrame) error {
	// 查找方法
	method, err := i.findMethod(def, methodName, methodDescriptor)
	if nil != err {
		return fmt.Errorf("failed to find method: %w", err)
	}

	// 解析访问标记
	flagMap := accflag.ParseAccFlags(method.AccessFlags)
	// 是native方法
	if _, ok := flagMap[accflag.Native]; ok {
		// 查本地方发表
		nativeFunc, argCount := i.miniJvm.NativeMethodTable.FindMethod(methodName, methodDescriptor)
		if nil == nativeFunc {
			// 该本地方法尚未被支持
			return fmt.Errorf("unsupported native method '%s'", method)
		}

		// 从操作数栈取出argCount个参数
		argCount += 1
		args := make([]interface{}, 0, argCount)
		for ix := 0; ix < argCount; ix++ {
			arg, _ := lastFrame.opStack.Pop()
			args = append(args, arg)
		}

		// 将jvm指针放到参数里,给native方法访问jvm的能力
		args[argCount - 1] = i.miniJvm

		// 因为出栈顺序跟实际参数顺序是相反的, 所以需要反转数组
		for ix := 0; ix < argCount / 2; ix++ {
			args[ix], args[argCount - 1 - ix] = args[argCount - 1 - ix], args[ix]
		}

		i.miniJvm.DebugPrintHistory = append(i.miniJvm.DebugPrintHistory, args[1:]...)

		// 调用go函数
		nativeFunc(args...)

		return nil
	}

	// 提取code属性
	codeAttr, err := i.findCodeAttr(method)
	if nil != err {
		return fmt.Errorf("failed to extract code attr: %w", err)
	}

	// 创建栈帧
	frame := newMethodStackFrame(int(codeAttr.MaxStack), int(codeAttr.MaxLocals))

	// 如果没有上层栈帧
	if nil == lastFrame {
		// 说明是main方法
		// todo 暂时不保存参数

	} else {
		// 传参
		// 判断是不是static方法
		var localVarStartIndexOffset int
		_, isStatic := flagMap[accflag.Static]
		if isStatic {
			// 如果是static方法, 则参数列表从本地变量表的0开始塞入
			localVarStartIndexOffset = 0

		} else {
			// 如果不是static方法, 则参数列表从本地变量表的1开始塞入
			localVarStartIndexOffset = 1

		}

		// 取出方法描述符
		descriptor := def.ConstPool[method.DescriptorIndex].(*class.Utf8InfoConst).String()
		// 解析描述符
		argDespList, _ := class.ParseMethodDescriptor(descriptor)
		// 临时保存参数列表
		argList := make([]interface{}, 0, len(argDespList))
		// 按参数数量出栈, 取出参数
		for _, arg := range argDespList {
			// 是int/char参数
			if "I" == arg || "C" == arg {
				// 从上一个栈帧中出栈, 保存到新栈帧的localVarTable中
				op, _ := lastFrame.opStack.PopInt()
				argList = append(argList, op)

				// frame.localVariablesTable[ix + localVarStartIndexOffset] = op

			} else {
				return fmt.Errorf("unsupported argument descriptor '%s' in '%s'", arg, descriptor)
			}
		}

		// 反转参数列表(因出栈顺序跟实际参数顺序相反)
		for ix := 0; ix < len(argList) / 2; ix++ {
			argList[ix], argList[len(argList) - 1 - ix] = argList[len(argList) - 1 - ix], argList[ix]
		}

		// 放入变量曹
		for ix, arg := range argList {
			frame.localVariablesTable[ix + localVarStartIndexOffset] = arg
		}

		if !isStatic {
			// 将this引用塞入0的位置
			obj, _ := lastFrame.opStack.PopReference()
			frame.localVariablesTable[0] = obj
		}
	}



	// 执行字节码
	return i.executeInFrame(def, codeAttr, frame, lastFrame)
}

func (i *InterpretedExecutionEngine) executeInFrame(def *class.DefFile, codeAttr *class.CodeAttr, frame *MethodStackFrame, lastFrame *MethodStackFrame) error {
	isWideStatus := false
	for {
		// 取出pc指向的字节码
		byteCode := codeAttr.Code[frame.pc]

		exitLoop := false

		// 执行
		switch byteCode {
		case bcode.Iconst0:
			// 将x压栈
			frame.opStack.Push(0)
		case bcode.Iconst1:
			frame.opStack.Push(1)
		case bcode.Iconst2:
			frame.opStack.Push(2)
		case bcode.Iconst3:
			frame.opStack.Push(3)
		case bcode.Iconst4:
			frame.opStack.Push(4)
		case bcode.Iconst5:
			frame.opStack.Push(5)

		case bcode.Iaload:
			// 将int型数组指定索引的值推送至栈顶
			// Operand Stack
			//..., arrayref, index →
			//..., value
			arrIndex, _ := frame.opStack.PopInt()
			arrRef, _ := frame.opStack.PopReference()
			frame.opStack.Push(arrRef.Array.Data[arrIndex])

		case bcode.Caload:
			// 将char型数组指定索引的值推送至栈顶
			// Operand Stack
			//..., arrayref, index →
			//..., value
			arrIndex, _ := frame.opStack.PopInt()
			arrRef, _ := frame.opStack.PopReference()
			frame.opStack.Push(arrRef.Array.Data[arrIndex])

		case bcode.Istore1:
			// 将栈顶int型数值存入第二个本地变量
			top, _ := frame.opStack.PopInt()
			frame.localVariablesTable[1] = top
		case bcode.Istore2:
			// 将栈顶int型数值存入第3个本地变量
			top, _ := frame.opStack.PopInt()
			frame.localVariablesTable[2] = top
		case bcode.Istore3:
			// 将栈顶int型数值存入第4个本地变量
			top, _ := frame.opStack.PopInt()
			frame.localVariablesTable[3] = top

		case bcode.Lstore1:
			// 将栈顶long型数值存入本地变量
			top, _ := frame.opStack.Pop()
			frame.localVariablesTable[1] = top

		case bcode.Iload:
			// Load int from local variable
			// ilaod index
			index := codeAttr.Code[frame.pc + 1]
			frame.pc++

			frame.opStack.Push(frame.localVariablesTable[index])
		case bcode.Iload0:
			// 将第1个slot中的值压栈
			frame.opStack.Push(frame.localVariablesTable[0])
		case bcode.Iload1:
			frame.opStack.Push(frame.localVariablesTable[1])
		case bcode.Iload2:
			frame.opStack.Push(frame.localVariablesTable[2])
		case bcode.Iload3:
			frame.opStack.Push(frame.localVariablesTable[3])

		case bcode.Aload0:
			// 将第一个引用类型本地变量推送至栈顶
			ref := frame.GetLocalTableObjectAt(0)
			frame.opStack.Push(ref)
		case bcode.Aload1:
			ref := frame.GetLocalTableObjectAt(1)
			frame.opStack.Push(ref)
		case bcode.Aload2:
			// 将第3个引用类型本地变量推送至栈顶
			ref := frame.GetLocalTableObjectAt(2)
			frame.opStack.Push(ref)

		case bcode.Istore:
			// istore index
			// ..., value →
			idx := codeAttr.Code[frame.pc + 1]
			frame.pc++

			val, _ := frame.opStack.Pop()
			frame.localVariablesTable[idx] = val

		case bcode.Astore0:
			// 将栈顶引用型数值存入本地变量
			ref, _ := frame.opStack.Pop()
			frame.localVariablesTable[0] = ref
		case bcode.Astore1:
			// 将栈顶引用型数值存入本地变量
			ref, _ := frame.opStack.Pop()
			frame.localVariablesTable[1] = ref
		case bcode.Astore2:
			ref, _ := frame.opStack.Pop()
			frame.localVariablesTable[2] = ref

		case bcode.Iastore:
			// 在int数组中存储元素
			// stack: arrayref, index, value →
			val, _ := frame.opStack.PopInt()
			arrIndex, _ := frame.opStack.PopInt()
			arrRef, _ := frame.opStack.PopReference()

			arrRef.Array.Data[arrIndex] = val

		case bcode.Castore:
			// Store into char array
			// stack: arrayref, index, value →
			val, _ := frame.opStack.Pop()
			arrIndex, _ := frame.opStack.PopInt()
			arrRef, _ := frame.opStack.PopReference()
			arrRef.Array.Data[arrIndex] = val

		case bcode.Dup:
			// 复制栈顶数值并将复制值压入栈顶
			top, _ := frame.opStack.GetTop()
			frame.opStack.Push(top)

		case bcode.Iadd:
			// 取出栈顶2元素，相加，入栈
			op1, _ := frame.opStack.PopInt()
			op2, _ := frame.opStack.PopInt()
			sum := op1 + op2
			frame.opStack.Push(sum)

		case bcode.Bipush:
			// 将单字节的常量值(-128~127)推送至栈顶
			num := int8(codeAttr.Code[frame.pc + 1])
			frame.opStack.Push(int(num))
			frame.pc++

		case bcode.Sipush:
			// 将一个短整型常量(-32768~32767)推送至栈顶
			twoByteNum := codeAttr.Code[frame.pc + 1 : frame.pc + 1 + 2]
			frame.pc += 2

			var op int16
			err := binary.Read(bytes.NewBuffer(twoByteNum), binary.BigEndian, &op)
			if nil != err {
				return fmt.Errorf("failed to read offset for sipush: %w", err)
			}

			frame.opStack.Push(int(op))

		case bcode.Ifle:
			// 当栈顶int型数值小于等于0时跳转
			err := i.bcodeIfCompZero(frame, codeAttr, func(op1 int, op2 int) bool {
				return op1 <= op2
			})

			if nil != err {
				return fmt.Errorf("failed to execute 'ifle': %w", err)
			}
		case bcode.Iflt:
			// 当栈顶int型数值小于0时跳转
			err := i.bcodeIfCompZero(frame, codeAttr, func(op1 int, op2 int) bool {
				return op1 < op2
			})

			if nil != err {
				return fmt.Errorf("failed to execute 'iflt': %w", err)
			}
		case bcode.Ifge:
			// >= 0
			err := i.bcodeIfCompZero(frame, codeAttr, func(op1 int, op2 int) bool {
				return op1 >= op2
			})

			if nil != err {
				return fmt.Errorf("failed to execute 'ifge': %w", err)
			}
		case bcode.Ifgt:
			// > 0
			err := i.bcodeIfCompZero(frame, codeAttr, func(op1 int, op2 int) bool {
				return op1 > op2
			})

			if nil != err {
				return fmt.Errorf("failed to execute 'ifgt': %w", err)
			}
		case bcode.Ifne:
			// != 0
			err := i.bcodeIfCompZero(frame, codeAttr, func(op1 int, op2 int) bool {
				return op1 != op2
			})

			if nil != err {
				return fmt.Errorf("failed to execute 'ifne': %w", err)
			}
		case bcode.Ifeq:
			// == 0
			err := i.bcodeIfCompZero(frame, codeAttr, func(op1 int, op2 int) bool {
				return op1 == op2
			})

			if nil != err {
				return fmt.Errorf("failed to execute 'ifeq': %w", err)
			}

		case bcode.Ificmpgt:
			// 比较栈顶两int型数值大小, 当结果大于0时跳转
			err := i.bcodeIfComp(frame, codeAttr, func(op1 int, op2 int) bool {
				return op2 - op1 > 0
			})
			if nil != err {
				return fmt.Errorf("failed to execute 'ificmpgt': %w", err)
			}
		case bcode.Ificmple:
			// 比较栈顶两int型数值大小, 当结果<=0时跳转
			err := i.bcodeIfComp(frame, codeAttr, func(op1 int, op2 int) bool {
				// fmt.Printf("%v compare %v\n", op1, op2)
				return op2 - op1 <= 0
			})
			if nil != err {
				return fmt.Errorf("failed to execute 'ificmple': %w", err)
			}
		case bcode.Ificmplt:
			// 比较栈顶两int型数值大小, 当结果小于0时跳转
			err := i.bcodeIfComp(frame, codeAttr, func(op1 int, op2 int) bool {
				return op2 - op1 < 0
			})
			if nil != err {
				return fmt.Errorf("failed to execute 'ificmplt': %w", err)
			}
		case bcode.Ificmpge:
			// 比较栈顶两int型数值大小, 当结果大于等于0时跳转
			err := i.bcodeIfComp(frame, codeAttr, func(op1 int, op2 int) bool {
				return op2 - op1 >= 0
			})
			if nil != err {
				return fmt.Errorf("failed to execute 'ificmpge': %w", err)
			}
		case bcode.Ificmpeq:
			// 比较栈顶两int型数值大小, 当结果等于0时跳转
			err := i.bcodeIfComp(frame, codeAttr, func(op1 int, op2 int) bool {
				return op2 - op1 == 0
			})
			if nil != err {
				return fmt.Errorf("failed to execute 'ificmpeq': %w", err)
			}
		case bcode.Ificmpne:
			// 比较栈顶两int型数值大小, 当结果!=0时跳转
			err := i.bcodeIfComp(frame, codeAttr, func(op1 int, op2 int) bool {
				return op2 != op1
			})
			if nil != err {
				return fmt.Errorf("failed to execute 'ificmpne': %w", err)
			}

		case bcode.Isub:
			// ..., value1, value2 →
			// The int result is value1 - value2. The result is pushed onto the operand stack.
			val2, _ := frame.opStack.PopInt()
			val1, _ := frame.opStack.PopInt()
			val := val1 - val2

			frame.opStack.Push(val)


		case bcode.Iinc:
			// 将第op1个slot的变量增加op2
			// iinc  byte constbyte
			if !isWideStatus {
				op1 := codeAttr.Code[frame.pc + 1]
				op2 := int8(codeAttr.Code[frame.pc + 2])
				frame.pc += 2

				frame.localVariablesTable[op1] = frame.GetLocalTableIntAt(int(op1)) + int(op2)

			} else {
				// wide iinc byte1 byte2 constbyte1 constbyte2
				twoByteNum := codeAttr.Code[frame.pc + 1 : frame.pc + 1 + 2]
				var localVarIndex uint16
				err := binary.Read(bytes.NewBuffer(twoByteNum), binary.BigEndian, &localVarIndex)
				if nil != err {
					return fmt.Errorf("failed to read local_var_index for iinc_w: %w", err)
				}


				twoByteNum = codeAttr.Code[frame.pc + 1 + 2 : frame.pc + 1 + 2 + 2]
				var num int16
				err = binary.Read(bytes.NewBuffer(twoByteNum), binary.BigEndian, &num)
				if nil != err {
					return fmt.Errorf("failed to read byte12 for iinc_w: %w", err)
				}

				frame.pc += 4

				newVal := frame.GetLocalTableIntAt(int(localVarIndex)) + int(num)
				frame.localVariablesTable[localVarIndex] = newVal

				isWideStatus = false
			}


		case bcode.New:
			// 创建一个对象, 并将其引用值压入栈顶
			twoByteNum := codeAttr.Code[frame.pc + 1 : frame.pc + 1 + 2]
			frame.pc += 2

			var classCpIndex uint16
			err := binary.Read(bytes.NewBuffer(twoByteNum), binary.BigEndian, &classCpIndex)
			if nil != err {
				return fmt.Errorf("failed to read class_cp_index for 'new': %w", err)
			}

			// 常量池中找出引用的class信息
			classCp := def.ConstPool[classCpIndex].(*class.ClassInfoConstInfo)
			// 目标class全名
			targetClassFullName := def.ConstPool[classCp.FullClassNameIndex].(*class.Utf8InfoConst).String()
			// 加载
			targetDefClass, err := i.miniJvm.MethodArea.LoadClass(targetClassFullName)
			if nil != err {
				return fmt.Errorf("failed to load class for '%s': %w", targetClassFullName, err)
			}
			// new
			obj, err := class.NewObject(targetDefClass, i.miniJvm.MethodArea)
			if nil != err {
				return fmt.Errorf("failed to new object for '%s': %w", targetClassFullName, err)
			}
			// 压栈
			frame.opStack.Push(obj)


		case bcode.Goto:
			// 跳转
			twoByteNum := codeAttr.Code[frame.pc + 1 : frame.pc + 1 + 2]
			var offset int16
			err := binary.Read(bytes.NewBuffer(twoByteNum), binary.BigEndian, &offset)
			if nil != err {
				return fmt.Errorf("failed to read pc offset for 'goto': %w", err)
			}

			frame.pc = frame.pc + int(offset) - 1

		case bcode.Invokestatic:
			// 调用静态方法
			err := i.invokeStatic(def, frame, codeAttr)
			if nil != err {
				return fmt.Errorf("failed to execute 'invokestatic': %w", err)
			}

		case bcode.Invokespecial:
			// 调用超类构建方法, 实例初始化方法, 私有方法
			err := i.invokeSpecial(def, frame, codeAttr)
			if nil != err {
				return fmt.Errorf("failed to execute 'invokespecial': %w", err)
			}

		case bcode.Invokevirtual:
			// public method
			err := i.invokeVirtual(def, frame, codeAttr)
			if nil != err {
				return fmt.Errorf("failed to execute 'invokevirtual': %w", err)
			}

		case bcode.Invokeinterface:
			// invokeinterface
			// indexbyte1
			// indexbyte2
			// count
			// 0
			err := i.invokeInterface(def, frame, codeAttr)
			if nil != err {
				return fmt.Errorf("failed to execute 'invokeinterface': %w", err)
			}

		case bcode.Getstatic:
			// format: getstatic byte1 byte2
			// Operand Stack
			// ..., →
			// ..., value
			err := i.bcodeGetStatic(def, frame, codeAttr)
			if nil != err {
				return fmt.Errorf("failed to execute 'getstatic': %w", err)
			}

		case bcode.Putstatic:
			// putstatic b1 b2
			// Operand Stack
			//..., value →
			//...
			err := i.bcodePutStatic(def, frame, codeAttr)
			if nil != err {
				return fmt.Errorf("failed to execute 'putstatic': %w", err)
			}


		case bcode.Putfield:
			// 对象字段赋值
			twoByteNum := codeAttr.Code[frame.pc + 1 : frame.pc + 1 + 2]
			frame.pc += 2

			var fieldRefCpIndex uint16
			err := binary.Read(bytes.NewBuffer(twoByteNum), binary.BigEndian, &fieldRefCpIndex)
			if nil != err {
				return fmt.Errorf("failed to read field_ref_cp_index: %w", err)
			}

			// 取出引用的字段
			fieldRef := def.ConstPool[fieldRefCpIndex].(*class.FieldRefConstInfo)
			// 取出字段名
			nameAndType := def.ConstPool[fieldRef.NameAndTypeIndex].(*class.NameAndTypeConst)
			fieldName := def.ConstPool[nameAndType.NameIndex].(*class.Utf8InfoConst).String()

			// 赋值
			val, _ := frame.opStack.PopInt()
			ref, _ := frame.opStack.PopReference()
			ref.Object.ObjectFields[fieldName].FieldValue = val
			// thisObj.ObjectFields[fieldName].FieldValue, _ = frame.opStack.PopInt()

		case bcode.GetField:
			// 获取指定对象的实例域, 并将其压入栈顶
			twoByteNum := codeAttr.Code[frame.pc + 1 : frame.pc + 1 + 2]
			frame.pc += 2

			var fieldRefCpIndex uint16
			err := binary.Read(bytes.NewBuffer(twoByteNum), binary.BigEndian, &fieldRefCpIndex)
			if nil != err {
				return fmt.Errorf("failed to read field_ref_cp_index: %w", err)
			}

			// 取出引用的字段
			fieldRef := def.ConstPool[fieldRefCpIndex].(*class.FieldRefConstInfo)
			// 取出字段名
			nameAndType := def.ConstPool[fieldRef.NameAndTypeIndex].(*class.NameAndTypeConst)
			fieldName := def.ConstPool[nameAndType.NameIndex].(*class.Utf8InfoConst).String()

			// 取出引用的对象
			targetObjRef, _ := frame.opStack.PopReference()

			// 读取
			val := targetObjRef.Object.ObjectFields[fieldName].FieldValue
			// 压栈
			frame.opStack.Push(val)

		case bcode.Newarray:
			// newarray type(byte)
			// 取出数组类型
			arrayType := codeAttr.Code[frame.pc + 1]
			frame.pc += 1

			// 栈顶元素为数组长度
			arrLen, _ := frame.opStack.PopInt()

			arrRef, err := class.NewArray(arrLen, arrayType)
			if nil != err {
				return fmt.Errorf("failed to execute 'newarray': %w", err)
			}

			// 数组引用入栈
			frame.opStack.Push(arrRef)

		case bcode.Athrow:
			err := i.bcodeAthrow(def, frame, codeAttr)
			if nil != err {
				return fmt.Errorf("failed to execute 'athrow': %w", err)
			}

		case bcode.Ireturn:
			// 当前栈出栈, 值压如上一个栈
			op, _ := frame.opStack.PopInt()
			lastFrame.opStack.Push(op)

			exitLoop = true

		case bcode.Return:
			// 返回
			exitLoop = true

		case bcode.Wide:
			// 加宽下一个字节码
			isWideStatus = true

		default:
			return fmt.Errorf("unsupported byte code %s", hex.EncodeToString([]byte{byteCode}))
		}

		if exitLoop {
			break
		}

		// 移动程序计数器
		frame.pc++
	}

	return nil
}

func (i *InterpretedExecutionEngine) invokeStatic(def *class.DefFile, frame *MethodStackFrame, codeAttr *class.CodeAttr) error {
	twoByteNum := codeAttr.Code[frame.pc + 1 : frame.pc + 1 + 2]
	frame.pc += 2

	var methodRefCpIndex uint16
	err := binary.Read(bytes.NewBuffer(twoByteNum), binary.BigEndian, &methodRefCpIndex)
	if nil != err {
		return fmt.Errorf("failed to read method_ref_cp_index: %w", err)
	}

	// 取出引用的方法
	methodRef := def.ConstPool[methodRefCpIndex].(*class.MethodRefConstInfo)
	// 取出方法名
	nameAndType := def.ConstPool[methodRef.NameAndTypeIndex].(*class.NameAndTypeConst)
	methodName := def.ConstPool[nameAndType.NameIndex].(*class.Utf8InfoConst).String()
	// 描述符
	descriptor := def.ConstPool[nameAndType.DescIndex].(*class.Utf8InfoConst).String()
	// 取出方法所在的class
	classRef := def.ConstPool[methodRef.ClassIndex].(*class.ClassInfoConstInfo)
	// 取出目标class全名
	targetClassFullName := def.ConstPool[classRef.FullClassNameIndex].(*class.Utf8InfoConst).String()
	// 加载
	targetDef, err := i.miniJvm.findDefClass(targetClassFullName)
	if nil != err {
		return fmt.Errorf("failed to load class for '%s': %w", targetClassFullName, err)
	}

	// 调用
	return i.ExecuteWithFrame(targetDef, methodName, descriptor, frame)
}

func (i *InterpretedExecutionEngine) invokeSpecial(def *class.DefFile, frame *MethodStackFrame, codeAttr *class.CodeAttr) error {
	twoByteNum := codeAttr.Code[frame.pc + 1 : frame.pc + 1 + 2]
	frame.pc += 2

	var methodRefCpIndex uint16
	err := binary.Read(bytes.NewBuffer(twoByteNum), binary.BigEndian, &methodRefCpIndex)
	if nil != err {
		return fmt.Errorf("failed to read method_ref_cp_index: %w", err)
	}


	// 取出引用的方法
	methodRef := def.ConstPool[methodRefCpIndex].(*class.MethodRefConstInfo)
	// 取出方法名
	nameAndType := def.ConstPool[methodRef.NameAndTypeIndex].(*class.NameAndTypeConst)
	methodName := def.ConstPool[nameAndType.NameIndex].(*class.Utf8InfoConst).String()
	// 描述符
	descriptor := def.ConstPool[nameAndType.DescIndex].(*class.Utf8InfoConst).String()
	// 取出方法所在的class
	classRef := def.ConstPool[methodRef.ClassIndex].(*class.ClassInfoConstInfo)
	// 取出目标class全名
	targetClassFullName := def.ConstPool[classRef.FullClassNameIndex].(*class.Utf8InfoConst).String()
	// 加载
	targetDef, err := i.miniJvm.findDefClass(targetClassFullName)
	if nil != err {
		return fmt.Errorf("failed to load class for '%s': %w", targetClassFullName, err)
	}

	if "<init>" == methodName {
		// 忽略构造器
		// 消耗一个引用
		frame.opStack.PopReference()
		return nil
	}

	// 调用
	return i.ExecuteWithFrame(targetDef, methodName, descriptor, frame)
}

func (i *InterpretedExecutionEngine) invokeVirtual(def *class.DefFile, frame *MethodStackFrame, codeAttr *class.CodeAttr) error {
	twoByteNum := codeAttr.Code[frame.pc + 1 : frame.pc + 1 + 2]
	frame.pc += 2

	var methodRefCpIndex uint16
	err := binary.Read(bytes.NewBuffer(twoByteNum), binary.BigEndian, &methodRefCpIndex)
	if nil != err {
		return fmt.Errorf("failed to read method_ref_cp_index: %w", err)
	}

	// 找到操作数栈中的第一个引用
	// 此引用即为实际类型
	targetObjRef, _ := frame.opStack.GetUntilObject()
	targetDef := targetObjRef.Object.DefFile

	// 取出引用的方法
	methodRef := def.ConstPool[methodRefCpIndex].(*class.MethodRefConstInfo)
	// 取出方法名
	nameAndType := def.ConstPool[methodRef.NameAndTypeIndex].(*class.NameAndTypeConst)
	methodName := def.ConstPool[nameAndType.NameIndex].(*class.Utf8InfoConst).String()
	// 描述符
	descriptor := def.ConstPool[nameAndType.DescIndex].(*class.Utf8InfoConst).String()


	//// 取出方法所在的class
	//classRef := def.ConstPool[methodRef.ClassIndex].(*class.ClassInfoConstInfo)
	//// 取出目标class全名
	//targetClassFullName := def.ConstPool[classRef.FullClassNameIndex].(*class.Utf8InfoConst).String()
	//// 加载
	//targetDef, err := i.miniJvm.findDefClass(targetClassFullName)
	//if nil != err {
	//	return fmt.Errorf("failed to load class for '%s': %w", targetClassFullName, err)
	//}

	// 取出栈顶对象引用
	// targetObj, _ := frame.opStack.PopReference()


	// 调用
	return i.ExecuteWithFrame(targetDef, methodName, descriptor, frame)
}

func (i *InterpretedExecutionEngine) invokeInterface(def *class.DefFile, frame *MethodStackFrame, codeAttr *class.CodeAttr) error {
	// invokeinterface
	// indexbyte1
	// indexbyte2
	// count
	// 0

	// 读取方法引用索引
	twoByteNum := codeAttr.Code[frame.pc + 1 : frame.pc + 1 + 2]
	var interfaceConstIndex int16
	err := binary.Read(bytes.NewBuffer(twoByteNum), binary.BigEndian, &interfaceConstIndex)
	if nil != err {
		return fmt.Errorf("failed to read interface_const_index for 'invokeinterface': %w", err)
	}

	// 多消耗2 byte
	twoByteNum = codeAttr.Code[frame.pc + 1 + 2 : frame.pc + 1 + 2 + 2]
	var nothing int16
	err = binary.Read(bytes.NewBuffer(twoByteNum), binary.BigEndian, &nothing)
	if nil != err {
		return fmt.Errorf("failed to read interface_const_index.nothing for 'invokeinterface': %w", err)
	}

	// 移动计数器
	frame.pc += 4

	// 取出接口方法引用
	interfaceMethodRef := def.ConstPool[interfaceConstIndex].(*class.InterfaceMethodConst)
	nameAndType := def.ConstPool[interfaceMethodRef.NameAndTypeIndex].(*class.NameAndTypeConst)

	targetMethodName := def.ConstPool[nameAndType.NameIndex].(*class.Utf8InfoConst).String()
	targetDescriptor := def.ConstPool[nameAndType.DescIndex].(*class.Utf8InfoConst).String()

	// 出栈取出对象引用
	ref, _ := frame.opStack.GetUntilObject()
	return i.ExecuteWithFrame(ref.Object.DefFile, targetMethodName, targetDescriptor, frame)
}

// 解释athrow指令
func (i *InterpretedExecutionEngine) bcodeAthrow(def *class.DefFile, frame *MethodStackFrame, codeAttr *class.CodeAttr) error {
	// 栈顶一定是异常对象引用
	ref, _ := frame.opStack.GetTopObject()

	// 栈顶异常全名
	thisExpInfo, _ := ref.Object.DefFile.ConstPool[ref.Object.DefFile.ThisClass].(*class.ClassInfoConstInfo)
	thisExpFullName := ref.Object.DefFile.ConstPool[thisExpInfo.FullClassNameIndex].(*class.Utf8InfoConst).String()

	// 查异常表
	for _, expTable := range codeAttr.ExceptionTable {
		// 确保当前pc是在范围内
		if frame.pc < int(expTable.StartPc) || frame.pc > int(expTable.EndPc) {
			continue
		}

		// 取出目标异常类型
		targetExpInfo := def.ConstPool[expTable.CatchType].(*class.ClassInfoConstInfo)
		// 目标异常全名
		targetExpFullName := def.ConstPool[targetExpInfo.FullClassNameIndex].(*class.Utf8InfoConst).String()

		// 判断跟栈顶异常是否匹配
		if targetExpFullName == thisExpFullName {
			// 修改pc实现跳转
			frame.pc = int(expTable.HandlerPc) - 1
			return nil
		}
	}

	return nil
}

// 读取static字段
// format: getstatic byte1 byte2
// Operand Stack
// ..., →
// ..., value
func (i *InterpretedExecutionEngine) bcodeGetStatic(def *class.DefFile, frame *MethodStackFrame, codeAttr *class.CodeAttr) error {
	// 静态字段在cp里的index
	twoByte := codeAttr.Code[frame.pc + 1 : frame.pc + 1 + 2]
	var fieldCpIndex int16
	err := binary.Read(bytes.NewBuffer(twoByte), binary.BigEndian, &fieldCpIndex)
	if nil != err {
		return fmt.Errorf("failed to read static field index: %w", err)
	}

	frame.pc += 2

	// 静态字段cp信息
	fieldInfo := def.ConstPool[fieldCpIndex].(*class.FieldRefConstInfo)
	// 取出字段所属class
	targetClassInfo := def.ConstPool[fieldInfo.ClassIndex].(*class.ClassInfoConstInfo)
	// 目标class全名
	targetClassFullName := def.ConstPool[targetClassInfo.FullClassNameIndex].(*class.Utf8InfoConst).String()
	// 加载
	targetClassDef, err := i.miniJvm.findDefClass(targetClassFullName)
	if nil != err {
		return fmt.Errorf("failed to load target class '%s':%w", targetClassFullName, err)
	}

	// 字段nameAndType
	nameAndTypeInfo := def.ConstPool[fieldInfo.NameAndTypeIndex].(*class.NameAndTypeConst)
	fieldName := def.ConstPool[nameAndTypeInfo.NameIndex].(*class.Utf8InfoConst).String()
	// fieldDesc := def.ConstPool[nameAndTypeInfo.DescIndex].(*class.Utf8InfoConst).String()

	// 查找目标字段
	objectField := targetClassDef.ParsedStaticFields[fieldName]
	// 压栈
	frame.opStack.Push(objectField)

	return nil
}

func (i *InterpretedExecutionEngine) bcodePutStatic(def *class.DefFile, frame *MethodStackFrame, codeAttr *class.CodeAttr) error {
	// 静态字段在cp里的index
	twoByte := codeAttr.Code[frame.pc + 1 : frame.pc + 1 + 2]
	var fieldCpIndex int16
	err := binary.Read(bytes.NewBuffer(twoByte), binary.BigEndian, &fieldCpIndex)
	if nil != err {
		return fmt.Errorf("failed to read static field index: %w", err)
	}

	frame.pc += 2

	// 静态字段cp信息
	fieldInfo := def.ConstPool[fieldCpIndex].(*class.FieldRefConstInfo)
	// 取出字段所属class
	targetClassInfo := def.ConstPool[fieldInfo.ClassIndex].(*class.ClassInfoConstInfo)
	// 目标class全名
	targetClassFullName := def.ConstPool[targetClassInfo.FullClassNameIndex].(*class.Utf8InfoConst).String()
	// 加载
	targetClassDef, err := i.miniJvm.findDefClass(targetClassFullName)
	if nil != err {
		return fmt.Errorf("failed to load target class '%s':%w", targetClassFullName, err)
	}

	// 字段nameAndType
	nameAndTypeInfo := def.ConstPool[fieldInfo.NameAndTypeIndex].(*class.NameAndTypeConst)
	fieldName := def.ConstPool[nameAndTypeInfo.NameIndex].(*class.Utf8InfoConst).String()
	// fieldDesc := def.ConstPool[nameAndTypeInfo.DescIndex].(*class.Utf8InfoConst).String()


	// 出栈
	val, _ := frame.opStack.Pop()

	// set字段
	targetClassDef.ParsedStaticFields[fieldName] = class.NewObjectField(val)

	return nil
}

func (i *InterpretedExecutionEngine) bcodeIfComp(frame *MethodStackFrame, codeAttr *class.CodeAttr, gotoJudgeFunc func(int, int) bool) error {
	// 比较栈顶两int型数值大小

	// 待比较的数
	x, _ := frame.opStack.PopInt()
	y, _ := frame.opStack.PopInt()

	// 跳转的偏移量
	twoByteNum := codeAttr.Code[frame.pc + 1 : frame.pc + 1 + 2]
	var offset int16
	err := binary.Read(bytes.NewBuffer(twoByteNum), binary.BigEndian, &offset)
	if nil != err {
		return fmt.Errorf("failed to read offset for if_icmpgt: %w", err)
	}

	if gotoJudgeFunc(x, y) {
		frame.pc = frame.pc + int(offset) - 1

	} else {
		frame.pc += 2
	}

	return nil
}

func (i *InterpretedExecutionEngine) bcodeIfCompZero(frame *MethodStackFrame, codeAttr *class.CodeAttr, gotoJudgeFunc func(int, int) bool) error {
	// 当栈顶int型数值小于0时跳转
	// 跳转的偏移量
	twoByteNum := codeAttr.Code[frame.pc + 1 : frame.pc + 1 + 2]
	var offset int16
	err := binary.Read(bytes.NewBuffer(twoByteNum), binary.BigEndian, &offset)
	if nil != err {
		return fmt.Errorf("failed to read offset for if_icmpgt: %w", err)
	}

	op, _ := frame.opStack.PopInt()
	if gotoJudgeFunc(op, 0) {
		frame.pc = frame.pc + int(offset) - 1

	} else {
		frame.pc += 2
	}

	return nil
}

func (i *InterpretedExecutionEngine) findCodeAttr(method *class.MethodInfo) (*class.CodeAttr, error) {
	for _, attrGeneric := range method.Attrs {
		attr, ok := attrGeneric.(*class.CodeAttr)
		if ok {
			return attr, nil
		}
	}

	// return nil, errors.New("no node attr in method")
	// native方法没有code属性
	return nil, nil
}

// 查找方法定义;
// def: 当前class定义
// methodName: 目标方法简单名
// methodDescriptor: 目标方法描述符
func (i *InterpretedExecutionEngine) findMethod(def *class.DefFile, methodName string, methodDescriptor string) (*class.MethodInfo, error) {
	currentClassDef := def
	for {
		for _, method := range currentClassDef.Methods {
			name := currentClassDef.ConstPool[method.NameIndex].(*class.Utf8InfoConst).String()
			descriptor := currentClassDef.ConstPool[method.DescriptorIndex].(*class.Utf8InfoConst).String()
			// 匹配简单名和描述符
			if name == methodName && descriptor == methodDescriptor {
				return method, nil
			}
		}

		// 从父类中寻找
		parentClassRef := currentClassDef.ConstPool[currentClassDef.SuperClass].(*class.ClassInfoConstInfo)
		// 取出父类全名
		targetClassFullName := currentClassDef.ConstPool[parentClassRef.FullClassNameIndex].(*class.Utf8InfoConst).String()
		if "java/lang/Object" == targetClassFullName {
			break
		}

		// 加载父类
		parentDef, err := i.miniJvm.findDefClass(targetClassFullName)
		if nil != err {
			return nil, fmt.Errorf("failed to load superclass '%s': %w", targetClassFullName, err)
		}

		currentClassDef = parentDef
	}


	return nil, fmt.Errorf("method '%s' not found", methodName)
}

func NewInterpretedExecutionEngine(vm *MiniJvm) *InterpretedExecutionEngine {
	return &InterpretedExecutionEngine{
		miniJvm:     vm,
		// methodStack: NewMethodStack(1024),
	}
}

