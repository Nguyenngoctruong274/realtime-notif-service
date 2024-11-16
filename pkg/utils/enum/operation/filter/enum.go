package filter

import (
	"fmt"
)

type Enum struct {
	Name string
	Data string
}

func None() Enum {
	return Enum{}
}

func OperationEQ() Enum {
	return Enum{
		Name: "eq",
		Data: "eq",
	}
}

func OperationNEQ() Enum {
	return Enum{
		Name: "neq",
		Data: "neq",
	}
}

func OperationLT() Enum {
	return Enum{
		Name: "lt",
		Data: "lt",
	}
}

func OperationLTE() Enum {
	return Enum{
		Name: "lte",
		Data: "lte",
	}
}

func OperationGT() Enum {
	return Enum{
		Name: "gt",
		Data: "gt",
	}
}

func OperationGTE() Enum {
	return Enum{
		Name: "gte",
		Data: "gte",
	}
}

func OperationIN() Enum {
	return Enum{
		Name: "in",
		Data: "in",
	}
}

func OperationLike() Enum {
	return Enum{
		Name: "like",
		Data: "like",
	}
}

func OperationILike() Enum {
	return Enum{
		Name: "ilike",
		Data: "ilike",
	}
}

func GetListOperator() []string {
	return []string{
		OperationEQ().Data,
		OperationNEQ().Data,
		OperationLT().Data,
		OperationLTE().Data,
		OperationGT().Data,
		OperationGTE().Data,
		OperationIN().Data,
		OperationLike().Data,
		OperationILike().Data,
	}
}

type FilterOperation struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

func ProcessOperator(op FilterOperation) (string, interface{}) {
	switch op.Operator {
	case OperationEQ().Data:
		return fmt.Sprintf("%s = ?", op.Field), op.Value
	case OperationNEQ().Data:
		return fmt.Sprintf("%s != ?", op.Field), op.Value
	case OperationLT().Data:
		return fmt.Sprintf("%s < ?", op.Field), op.Value
	case OperationLTE().Data:
		return fmt.Sprintf("%s <= ?", op.Field), op.Value
	case OperationGT().Data:
		return fmt.Sprintf("%s > ?", op.Field), op.Value
	case OperationGTE().Data:
		return fmt.Sprintf("%s >= ?", op.Field), op.Value
	case OperationIN().Data:
		return fmt.Sprintf("%s IN (?)", op.Field), op.Value
	case OperationLike().Data:
		return fmt.Sprintf("%s like ?", op.Field), fmt.Sprintf("%%%v%%", op.Value)
	case OperationILike().Data:
		return fmt.Sprintf("%s ilike ?", op.Field), fmt.Sprintf("%%%v%%", op.Value)
	default:
		// Xử lý toán tử mặc định hoặc trả về một giá trị tùy ý
		return fmt.Sprintf("%s = ?", op.Field), op.Value
	}
}
