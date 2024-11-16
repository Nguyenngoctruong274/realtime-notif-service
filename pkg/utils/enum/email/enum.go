package email

type Enum struct {
	Name string
	Data string
}

func None() Enum {
	return Enum{}
}

func EmailYes4All() Enum {
	return Enum{
		Name: "@yes4all.com",
		Data: "@yes4all.com",
	}
}

func EmailVitoxVn() Enum {
	return Enum{
		Name: "@vitoxvn.com",
		Data: "@vitoxvn.com",
	}
}

func GetListEmail() []string {
	return []string{
		EmailYes4All().Data,
		EmailVitoxVn().Data,
	}
}

func GetEnumByName(data string) Enum {
	list := map[string]Enum{
		EmailYes4All().Data: EmailYes4All(),
		EmailVitoxVn().Data: EmailVitoxVn(),
	}
	if value, found := list[data]; found {
		return value
	}
	return None()
}
