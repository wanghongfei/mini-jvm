package utils

func InterfaceArrayToRuneArray(from []interface{}) []rune {
	runeArr := make([]rune, len(from))
	for ix, c := range from {
		runeArr[ix] = c.(rune)
	}

	return runeArr
}

func FillInterfaceArrayRune(from []interface{}, to []rune) {
	for ix, r := range to {
		from[ix] = r
	}
}
