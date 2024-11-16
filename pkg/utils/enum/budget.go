// nolint
package enum

type BudgetType struct {
	ID   int    `json:"id"`
	Path string `json:"path"`
	Name string `json:"name"`
}

func SponsoredProducts() BudgetType {
	return BudgetType{
		ID:   1,
		Path: "sp",
		Name: "Sponsored Products",
	}
}

func SponsoredBrands() BudgetType {
	return BudgetType{
		ID:   2,
		Path: "sb",
		Name: "Sponsored Brands",
	}
}

func SponsoredDisplay() BudgetType {
	return BudgetType{
		ID:   3,
		Path: "sd",
		Name: "Sponsored Display",
	}
}

func BudgetTypeIDByPath(path string) int {
	switch path {
	case SponsoredBrands().Path:
		return SponsoredBrands().ID
	case SponsoredDisplay().Path:
		return SponsoredDisplay().ID
	default:
		return SponsoredProducts().ID
	}
}
