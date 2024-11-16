package enum_notification

type Enum struct {
	Name string
	Data string
}

func None() Enum {
	return Enum{}
}

func UpdateAction() Enum {
	return Enum{
		Name: "update",
		Data: "update",
	}
}

func ErrorAction() Enum {
	return Enum{
		Name: "error",
		Data: "error",
	}
}

func AddAction() Enum {
	return Enum{
		Name: "add",
		Data: "add",
	}
}

func ListAction() Enum {
	return Enum{
		Name: "list",
		Data: "list",
	}
}

func GetEnumByData(data string) Enum {
	list := map[string]Enum{
		ListAction().Data:   ListAction(),
		AddAction().Data:    AddAction(),
		UpdateAction().Data: UpdateAction(),
	}
	if value, found := list[data]; found {
		return value
	}
	return None()
}
