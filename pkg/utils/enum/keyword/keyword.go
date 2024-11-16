package keyword

type Enum struct {
	Name string
	Data string
}

func None() Enum {
	return Enum{}
}

func MatchTypeTheme() Enum {
	return Enum{
		Name: "theme",
		Data: "theme",
	}
}

func ThemeTypeLandingPage() Enum {
	return Enum{
		Name: "keywords related to your landing pages",
		Data: "KEYWORDS_RELATED_TO_YOUR_LANDING_PAGES",
	}
}

func ThemeTypeBrand() Enum {
	return Enum{
		Name: "keywords related to your brand",
		Data: "KEYWORDS_RELATED_TO_YOUR_BRAND",
	}
}

func GetListMatchTypeTheme(data string) Enum {
	list := map[string]Enum{
		MatchTypeTheme().Data: MatchTypeTheme(),
	}

	if value, found := list[data]; found {
		return value
	}
	return None()
}

func GetListTypeTheme(data string) Enum {
	list := map[string]Enum{
		ThemeTypeLandingPage().Data: ThemeTypeLandingPage(),
		ThemeTypeBrand().Data:       ThemeTypeBrand(),
	}

	if value, found := list[data]; found {
		return value
	}
	return None()
}
