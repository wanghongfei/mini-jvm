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
	return i.execute(def, methodName, "([Ljava/lang/String;)V", nil)
}

func (i *InterpretedExecutionEngine) execute(def *class.DefFile, methodName string, methodDescriptor string, lastFrame *MethodStackFrame) error {
	// 查找方法
	method, err := i.findMethod(def, methodName, methodDescriptor)
	if nil != err {
		return fmt.Errorf("failed to find method: %w", err)
	}

	// 解析访问标记
	flagMap := accflag.ParseAccFlags(method.AccessFlags)
	if _, ok := flagMap[accflag.Native]; ok {
		// 特殊处理输出函数, 因为System.out太复杂了
		if "print" == methodName {
			data, _ := lastFrame.opStack.PopInt()
			fmt.Println(data)
			return nil
		}

		return fmt.Errorf("native method '%s' is unsupported", methodName)
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
		argList, _ := class.ParseMethodDescriptor(descriptor)
		// 将参数保存到新栈帧的本地变量表中
		for ix, arg := range argList {
			// 是int参数
			if "I" == arg {
				// 从上一个栈帧中出栈, 保存到新栈帧的localVarTable中
				op, _ := lastFrame.opStack.PopInt()
				frame.localVariablesTable[ix + localVarStartIndexOffset] = op

			} else {
				return fmt.Errorf("unsupported argument descriptor '%s' in '%s'", arg, descriptor)
			}
		}

		if !isStatic {
			// 将this引用塞入0的位置
			obj, _ := lastFrame.opStack.PopObject()
			frame.localVariablesTable[0] = obj
		}
	}



	// 执行字节码
	return i.executeInFrame(def, codeAttr, frame, lastFrame)
}

func (i *InterpretedExecutionEngine) executeInFrame(def *class.DefFile, codeAttr *class.CodeAttr, frame *MethodStackFrame, lastFrame *MethodStackFrame) error {
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
			num := codeAttr.Code[frame.pc + 1]
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

		case bcode.Ificmpgt:
			// 比较栈顶两int型数值大小, 当结果大于0时跳转

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

			if y - x > 0 {
				frame.pc = frame.pc + int(offset) - 1

			} else {
				frame.pc += 2
			}

		case bcode.Ificmple:
			// 比较栈顶两int型数值大小, 当结果<=0时跳转

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

			if y - x <= 0 {
				frame.pc = frame.pc + int(offset) - 1

			} else {
				frame.pc += 2
			}

		case bcode.Iinc:
			// 将第op1个slot的变量增加op2
			op1 := codeAttr.Code[frame.pc + 1]
			op2 := codeAttr.Code[frame.pc + 2]
			frame.pc += 2

			// frame.localVariablesTable[op1] = frame.localVariablesTable[op1] + int(op2)
			frame.localVariablesTable[op1] = frame.GetLocalTableIntAt(int(op1)) + int(op2)

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
			obj, err := class.NewObject(targetDefClass)
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
			obj, _ := frame.opStack.PopObject()
			obj.ObjectFields[fieldName].FieldValue = val
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
			targetObj, _ := frame.opStack.PopObject()

			// 读取
			val := targetObj.ObjectFields[fieldName].FieldValue
			// 压栈
			frame.opStack.Push(val)


		case bcode.Ireturn:
			// 当前栈出栈, 值压如上一个栈
			op, _ := frame.opStack.PopInt()
			lastFrame.opStack.Push(op)

			exitLoop = true

		case bcode.Return:
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
	return i.execute(targetDef, methodName, descriptor, frame)
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
		frame.opStack.PopObject()
		return nil
	}

	// 调用
	return i.execute(targetDef, methodName, descriptor, frame)
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
	// 次引用即为实际类型
	targetObj, _ := frame.opStack.GetUntilObject()
	targetDef := targetObj.DefFile

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
	// targetObj, _ := frame.opStack.PopObject()


	// 调用
	return i.execute(targetDef, methodName, descriptor, frame)
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
			name := def.ConstPool[method.NameIndex].(*class.Utf8InfoConst).String()
			descriptor := def.ConstPool[method.DescriptorIndex].(*class.Utf8InfoConst).String()
			// 匹配简单名和描述符
			if name == methodName && descriptor == methodDescriptor {
				return method, nil
			}
		}

		// 从父类中寻找
		parentClassRef := def.ConstPool[def.SuperClass].(*class.ClassInfoConstInfo)
		// 取出父类全名
		targetClassFullName := def.ConstPool[parentClassRef.FullClassNameIndex].(*class.Utf8InfoConst).String()
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

