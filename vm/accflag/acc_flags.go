package accflag

const (
	Public = 0x0001
	Private = 0x0002
	Protected = 0x0004
	Static = 0x0008
	Final = 0x0010
	Synchronized = 0x0020
	Bridge = 0x0040
	Varargs = 0x0080
	Native = 0x0100
	Abstarct = 0x0400
	Strict = 0x0800
	Synthetic = 0x1000
)

// 解析访问标记;
// return: key[标记值]任意值
func ParseAccFlags(flagBits uint16) map[int]interface{} {
	flagMap := make(map[int]interface{})

	if flagBits & Public > 0 {
		flagMap[Public] = struct {}{}
	}
	if flagBits & Private > 0 {
		flagMap[Private] = struct {}{}
	}
	if flagBits & Protected > 0 {
		flagMap[Protected] = struct {}{}
	}
	if flagBits & Static > 0 {
		flagMap[Static] = struct {}{}
	}
	if flagBits & Final > 0 {
		flagMap[Final] = struct {}{}
	}
	if flagBits & Synchronized > 0 {
		flagMap[Synchronized] = struct {}{}
	}
	if flagBits & Bridge > 0 {
		flagMap[Bridge] = struct {}{}
	}
	if flagBits & Varargs > 0 {
		flagMap[Varargs] = struct {}{}
	}
	if flagBits & Native > 0 {
		flagMap[Native] = struct {}{}
	}
	if flagBits & Abstarct > 0 {
		flagMap[Abstarct] = struct {}{}
	}
	if flagBits & Strict > 0 {
		flagMap[Strict] = struct {}{}
	}
	if flagBits & Synthetic > 0 {
		flagMap[Synthetic] = struct {}{}
	}

	return flagMap
}
