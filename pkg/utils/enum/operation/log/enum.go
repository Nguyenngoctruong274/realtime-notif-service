package log

type Enum struct {
	Name string
	Data string
}

func None() Enum {
	return Enum{}
}

func OperationPortfolio() Enum {
	return Enum{
		Name: "Portfolio",
		Data: "Portfolio",
	}
}

func OperationPortfolioUser() Enum {
	return Enum{
		Name: "PortfolioUser",
		Data: "PortfolioUser",
	}
}

func ActionPortfolioCreate() Enum {
	return Enum{
		Name: "Portfolio:Create",
		Data: "Portfolio:Create",
	}
}

func ActionPortfolioUpdate() Enum {
	return Enum{
		Name: "Portfolio:Update",
		Data: "Portfolio:Update",
	}
}

func ActionPortfolioDelete() Enum {
	return Enum{
		Name: "Portfolio:Delete",
		Data: "Portfolio:Delete",
	}
}

func ActionPortfolioDisable() Enum {
	return Enum{
		Name: "Portfolio:Disable",
		Data: "Portfolio:Disable",
	}
}

func ActionPortfolioEnable() Enum {
	return Enum{
		Name: "Portfolio:Enable",
		Data: "Portfolio:Enable",
	}
}

func ActionPortfolioUserCreate() Enum {
	return Enum{
		Name: "PortfolioUser:Create",
		Data: "PortfolioUser:Create",
	}
}

func ActionPortfolioUserDelete() Enum {
	return Enum{
		Name: "PortfolioUser:Delete",
		Data: "PortfolioUser:Delete",
	}
}

func ActionNotificationIncidentsCreate() Enum {
	return Enum{
		Name: "NotificationIncidents:Create",
		Data: "NotificationIncidents:Create",
	}
}

func ActionNotificationIncidentsUpdate() Enum {
	return Enum{
		Name: "NotificationIncidents:Update",
		Data: "NotificationIncidents:Update",
	}
}

func NotificationIncidents() Enum {
	return Enum{
		Name: "NotificationIncidents",
		Data: "NotificationIncidents",
	}
}

func GetEnumByName(data string) Enum {
	list := map[string]Enum{
		ActionPortfolioCreate().Data:     ActionPortfolioCreate(),
		ActionPortfolioUpdate().Data:     ActionPortfolioUpdate(),
		ActionPortfolioDelete().Data:     ActionPortfolioDelete(),
		ActionPortfolioDisable().Data:    ActionPortfolioDisable(),
		ActionPortfolioEnable().Data:     ActionPortfolioEnable(),
		ActionPortfolioUserCreate().Data: ActionPortfolioUserCreate(),
		ActionPortfolioUserDelete().Data: ActionPortfolioUserDelete(),
	}
	if value, found := list[data]; found {
		return value
	}
	return None()
}
