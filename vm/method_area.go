package vm

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/wanghongfei/mini-jvm/utils"
	"github.com/wanghongfei/mini-jvm/vm/accflag"
	"github.com/wanghongfei/mini-jvm/vm/class"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

var ClassIgnoredErr = errors.New("ignored")

// 方法区
type MethodArea struct {
	Jvm *MiniJvm

	// 类路径
	ClassPaths []string

	// key: 类的选限定性名
	// val: 加载完成后的DefFile
	// 因为有可能在其他goroutine中加载类, 所以需要加锁
	ClassMap map[string]*class.DefFile
	ClassMapLock sync.RWMutex

	// 忽略的class的全名, 遇到这些class时不触发加载逻辑
	IgnoredClasses map[string]interface{}
}

func NewMethodArea(jvm *MiniJvm, classpaths []string, ignoredClasses []string) (*MethodArea, error) {
	if nil == classpaths || len(classpaths) == 0 {
		return nil, fmt.Errorf("invalid classpath: %v", classpaths)
	}

	res := &MethodArea{
		Jvm: jvm,
		ClassPaths: classpaths,
		ClassMap: make(map[string]*class.DefFile),
		IgnoredClasses: make(map[string]interface{}),
	}

	if nil != ignoredClasses {
		for _, name := range ignoredClasses {
			res.IgnoredClasses[name] = struct {}{}
		}
	}

	return res, nil
}

// 从classpath中加载一个类
// fullname: 全限定性名
func (m *MethodArea) LoadClass(fullyQualifiedName string) (*class.DefFile, error) {
	utils.LogInfoPrintf("load class: %s", fullyQualifiedName)

	// 查忽略列表
	if _, ok := m.IgnoredClasses[fullyQualifiedName]; ok {
		// 此class被忽略
		return nil, ClassIgnoredErr
	}

	// 先从已加载的类中寻找
	m.ClassMapLock.RLock()
	targetClassDef, ok := m.ClassMap[fullyQualifiedName]
	m.ClassMapLock.RUnlock()
	if ok {
		utils.LogInfoPrintf("load class from cache: %s", fullyQualifiedName)
		return targetClassDef, nil
	}

	var defFile *class.DefFile

	// 从classpath寻找
	filepath, err := m.findClassFilePath(fullyQualifiedName)
	if nil != err {
		// 没找到
		// 从jar中寻找
		classBuf, err := m.findClassBuf(fullyQualifiedName)
		if nil != err {
			// 还没找到
			return nil, err
		}

		// 找到了
		// 加载class
		defFile, err = class.LoadClassBuf(classBuf)
		if nil != err {
			return nil, fmt.Errorf("unabled to load class %s: %w", fullyQualifiedName, err)
		}

	} else {
		defFile, err = class.LoadClassFile(filepath)
		if nil != err {
			return nil, fmt.Errorf("unabled to load class %s: %w", fullyQualifiedName, err)
		}
	}

	m.ClassMapLock.Lock()
	m.ClassMap[fullyQualifiedName] = defFile
	m.ClassMapLock.Unlock()

	// 执行<clinit>方法
	err = m.Jvm.ExecutionEngine.ExecuteWithDescriptor(defFile, "<clinit>", "()V")
	if nil != err && "failed to find method" == err.Error() {
		return nil, fmt.Errorf("failed to execute <clinit> for class '%s':%w", fullyQualifiedName, err)
	}

	// 初始化虚方法表
	err = m.initVTable(defFile)
	if nil != err {
		return nil, fmt.Errorf("failed to init vtable for class '%s':%w", fullyQualifiedName, err)
	}

	return defFile, nil
}

func (m *MethodArea) findClassFilePath(fullyQualifiedName string) (string, error) {

	for _, cp := range m.ClassPaths {
		possiblePath := cp + "/" + fullyQualifiedName + ".class"
		_, err := os.Stat(possiblePath)
		if nil == err {
			return possiblePath, nil
		}

	}

	return "", fmt.Errorf("cannot found class '%s' in classpath", fullyQualifiedName)
}

func (m *MethodArea) findClassBuf(fullyQualifiedName string) ([]byte, error) {
	var classFileBuf []byte

	for _, cp := range m.ClassPaths {
		if !strings.HasSuffix(cp, ".jar") {
			continue
		}

		destName := fullyQualifiedName + ".class"

		// 构造访问zip文件所需要的函数
		predicate := func(f *zip.File) bool {
			return f.Name == destName
		}

		visitor := func(reader io.Reader) (bool, error) {
			buf, err := ioutil.ReadAll(reader)
			classFileBuf = buf

			return true, err
		}

		utils.VisitZip(cp, predicate, visitor)
	}

	if 0 != len(classFileBuf) {
		return classFileBuf, nil
	}

	return nil, fmt.Errorf("cannot found class '%s' in classpath", fullyQualifiedName)
}

// 为指定class初始化虚方法表;
// 此方法同时也会递归触发父类虚方法表的初始化工作, 但不会重复初始化
func (m *MethodArea) initVTable(def *class.DefFile) error {
	def.VTable = make([]*class.VTableItem, 0, 5)

	// 取出父类引用信息
	superClassIndex := def.SuperClass
	// 没有父类
	if 0 == superClassIndex {
		// 遍历方法元数据, 添加到虚方法表中
		for _, methodInfo := range def.Methods {
			// 取出方法访问标记
			flagMap := accflag.ParseAccFlags(methodInfo.AccessFlags)
			_, isPublic := flagMap[accflag.Public]
			_, isProtected := flagMap[accflag.Protected]
			_, isNative := flagMap[accflag.Native]

			// 只添加public, protected, native方法
			if !isPublic && !isProtected && !isNative {
				// 跳过
				continue
			}

			// 取出方法名和描述符
			name := def.ConstPool[methodInfo.NameIndex].(*class.Utf8InfoConst).String()
			descriptor := def.ConstPool[methodInfo.DescriptorIndex].(*class.Utf8InfoConst).String()
			// 忽略构造方法
			if name == "<init>" {
				continue
			}


			newItem := &class.VTableItem{
				MethodName:       name,
				MethodDescriptor: descriptor,
				MethodInfo:       methodInfo,
			}
			def.VTable = append(def.VTable, newItem)
		}

		return nil
	}

	superClassInfo := def.ConstPool[superClassIndex].(*class.ClassInfoConstInfo)
	// 取出父类全名
	superClassFullName := def.ConstPool[superClassInfo.FullClassNameIndex].(*class.Utf8InfoConst).String()
	// 加载父类
	superDef, err := m.LoadClass(superClassFullName)
	if nil != err {
		return fmt.Errorf("cannot load parent class '%s'", superClassFullName)
	}

	// 判断父类虚方法表是否已经初始化过了
	if len(superDef.VTable) == 0 {
		// 没有初始化过
		// 初始化父类的虚方法表
		err = m.initVTable(superDef)
		if nil != err {
			return fmt.Errorf("cannot init vtable for parent class '%s':%w", superClassFullName, err)
		}
	}

	// 从父类虚方法表中继承元素
	for _, superItem := range superDef.VTable {
		subItem := &class.VTableItem{
			MethodName:       superItem.MethodName,
			MethodDescriptor: superItem.MethodDescriptor,
			MethodInfo:       superItem.MethodInfo,
		}

		def.VTable = append(def.VTable, subItem)
	}

	// 遍历自己的方法元数据, 替换或者追加虚方法表
	for _, methodInfo := range def.Methods {
		// 取出方法名和描述符
		name := def.ConstPool[methodInfo.NameIndex].(*class.Utf8InfoConst).String()
		descriptor := def.ConstPool[methodInfo.DescriptorIndex].(*class.Utf8InfoConst).String()
		// 忽略构造方法
		if name == "<init>" {
			continue
		}

		// 取出方法描述符
		flagMap := accflag.ParseAccFlags(methodInfo.AccessFlags)
		_, isPublic := flagMap[accflag.Public]
		_, isProtected := flagMap[accflag.Protected]
		_, isNative := flagMap[accflag.Native]
		// 只添加public, protected, native方法
		if !isPublic && !isProtected && !isNative {
			// 跳过
			continue
		}

		// 查找虚方法表中是否已经存在
		found := false
		for _, item := range def.VTable {
			if item.MethodName == name && item.MethodDescriptor == descriptor {
				// 说明def类重写了父类方法
				// 替换虚方法表当前项
				item.MethodInfo = methodInfo
				found = true
				break
			}
		}

		if !found {
			// 从父类继承的虚方法表中没找到此方法, 说明是子类的新方法, 追加
			newItem := &class.VTableItem{
				MethodName:       name,
				MethodDescriptor: descriptor,
				MethodInfo:       methodInfo,
			}
			def.VTable = append(def.VTable, newItem)
		}
	}

	return nil

}
