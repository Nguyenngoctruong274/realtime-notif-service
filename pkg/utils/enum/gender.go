//nolint:gomnd
package enum

type GenderEnum struct {
	Key   int
	Value string
	Name  string
}

func Male() GenderEnum {
	return GenderEnum{
		Key:   1,
		Value: "M",
		Name:  "anh",
	}
}

func Female() GenderEnum {
	return GenderEnum{
		Key:   2,
		Value: "F",
		Name:  "chị",
	}
}

func Undefined() GenderEnum {
	return GenderEnum{
		Key:   1000,
		Value: "O",
		Name:  "anh/chị",
	}
}

var keyToEnumLoop = map[int]GenderEnum{
	Male().Key:   Male(),
	Female().Key: Female(),
}

func GetGenderStringBykey(key int) (enum string) {
	switch key {
	case Male().Key:
		return Male().Name
	case Female().Key:
		return Female().Name
	default:
		return Undefined().Name
	}
}

func GetGenderByKey(key int) GenderEnum {
	result, ok := keyToEnumLoop[key]
	if ok {
		return result
	}

	return Undefined()
}
