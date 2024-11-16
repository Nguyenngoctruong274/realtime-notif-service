// nolint: gofmt,goimports
package number

import (
	"math"

	"github.com/shopspring/decimal"
)

func AddFloat64(current, add float64) (res float64) {
	currentDecimal := decimal.NewFromFloat(current)
	addDecimal := decimal.NewFromFloat(add)
	return currentDecimal.Add(addDecimal).InexactFloat64()
}

func ToFloat64CeilAfterDot(current float64, numberAfterDot float64) (res float64) {
	number := current * math.Pow(10, numberAfterDot)
	numberCeil := math.Ceil(number)
	res = numberCeil / math.Pow(10, numberAfterDot)
	return
}

func ToFloat64RoundAfterDot(current float64, numberAfterDot float64) (res float64) {
	number := current * math.Pow(10, numberAfterDot)
	numberCeil := math.Round(number)
	res = numberCeil / math.Pow(10, numberAfterDot)
	return
}
