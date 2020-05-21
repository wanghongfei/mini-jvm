package class

type Object struct {
	// class定义
	DefFile *DefFile

	// 实例数据
	ObjectFields []*ObjectField
}


type ObjectField struct {
	// 实例值
	FieldValue interface{}
	FieldType string
}
