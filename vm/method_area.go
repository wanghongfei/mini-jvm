package vm

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/wanghongfei/mini-jvm/utils"
	"github.com/wanghongfei/mini-jvm/vm/class"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

var ClassIgnoredErr = errors.New("ignored")

// 方法区
type MethodArea struct {
	Jvm *MiniJvm

	// 类路径
	ClassPaths []string

	// key: 类的选限定性名
	// val: 加载完成后的DefFile
	// todo 需要加锁
	ClassMap map[string]*class.DefFile

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
	// 查忽略列表
	if _, ok := m.IgnoredClasses[fullyQualifiedName]; ok {
		// 此class被忽略
		return nil, ClassIgnoredErr
	}

	// 先从已加载的类中寻找
	targetClassDef, ok := m.ClassMap[fullyQualifiedName]
	if ok {
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

	m.ClassMap[fullyQualifiedName] = defFile

	// 执行<clinit>方法
	err = m.Jvm.ExecutionEngine.ExecuteWithDescriptor(defFile, "<clinit>", "()V")
	if nil != err && "failed to find method" == err.Error() {
		return nil, fmt.Errorf("failed to execute <clinit> for class '%s':%w", fullyQualifiedName, err)
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
