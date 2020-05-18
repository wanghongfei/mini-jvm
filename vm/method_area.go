package vm

import (
	"fmt"
	"github.com/wanghongfei/mini-jvm/vm/class"
	"os"
)

// 方法区
type MethodArea struct {
	// 类路径
	ClassPaths []string

	// key: 类的选限定性名
	// val: 加载完成后的DefFile
	ClassMap map[string]*class.DefFile
}

func NewMethodArea(classpaths []string) (*MethodArea, error) {
	if nil == classpaths || len(classpaths) == 0 {
		return nil, fmt.Errorf("invalid classpath: %v", classpaths)
	}

	return &MethodArea{
		ClassPaths: classpaths,
		ClassMap: make(map[string]*class.DefFile),
	}, nil
}

// 从classpath中加载一个类
// fullname: 全限定性名
func (m *MethodArea) LoadClass(fullyQualifiedName string) error {
	filepath, err := m.findClassFile(fullyQualifiedName)
	if nil != err {
		return err
	}

	defFile, err := class.LoadClassFile(filepath)
	if nil != err {
		return fmt.Errorf("unabled to load class %s: %w", fullyQualifiedName, err)
	}

	m.ClassMap[fullyQualifiedName] = defFile
	return nil
}

func (m *MethodArea) findClassFile(fullyQualifiedName string) (string, error) {
	for _, cp := range m.ClassPaths {
		possiblePath := cp + "/" + fullyQualifiedName + ".class"
		_, err := os.Stat(possiblePath)
		if nil == err {
			return possiblePath, nil
		}
	}

	return "", fmt.Errorf("cannot found class '%s' in classpath", fullyQualifiedName)
}
