package portfolio

type Enum struct {
	Name   string
	Data   string
	DataDB string
}

func None() Enum {
	return Enum{}
}

func PortfolioStateEnabled() Enum {
	return Enum{
		Name: "enabled",
		Data: "enabled",
	}
}

func PortfolioStateDisabled() Enum {
	return Enum{
		Name: "disabled",
		Data: "disabled",
	}
}

func PortfolioStatePaused() Enum {
	return Enum{
		Name: "paused",
		Data: "paused",
	}
}

func PortfolioStateArchived() Enum {
	return Enum{
		Name: "archived",
		Data: "archived",
	}
}

func GetPortfolioStates() []string {
	return []string{
		PortfolioStateEnabled().Data,
	}
}
func PortfolioEmptyID() Enum {
	return Enum{
		Name: "1",
		Data: "1",
	}
}

func PortfolioEmptyIDEdit() Enum {
	return Enum{
		Name: "-1",
		Data: "-1",
	}
}

func PortfolioEmptyIDOperation() Enum {
	return Enum{
		Name:   "1",
		Data:   "1",
		DataDB: "",
	}
}

func PortfolioEmptyIDView() Enum {
	return Enum{
		Name: "",
		Data: "",
	}
}
func GetCurrencyCodes() []string {
	return []string{
		"USD",
		"CAD",
		"MXN",
		"BRL",
		"GBP",
		"JPY",
		"EUR",
		"AUD",
		"AED",
		"SEK",
		"PLN",
		"SGD",
		"TRY",
	}
}

func PortfolioPolicyMonthlyRecurring() Enum {
	return Enum{
		Name: "MonthlyRecurring",
		Data: "monthlyRecurring",
	}
}

func PortfolioPolicyDateRange() Enum {
	return Enum{
		Name: "DateRange",
		Data: "dateRange",
	}
}

func GetPolicies() []string {
	return []string{
		"dateRange",
		"monthlyRecurring",
	}
}

func GetEnumByName(data string) Enum {
	list := map[string]Enum{
		PortfolioStateEnabled().Data: PortfolioStateEnabled(),
	}
	if value, found := list[data]; found {
		return value
	}
	return None()
}
