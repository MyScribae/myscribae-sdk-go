package utilities

import (
	"fmt"
	"regexp"
	"strconv"
)

// MoneyValue is a type that represents a value in cents
type MoneyValue uint

// NewMoneyValue creates a new MoneyValue from a float64
func NewDollarValue(value float64) MoneyValue {
	return MoneyValue(value * 100)
}

// NewMoneyValue creates a new MoneyValue from a float64
func NewCentValue(value uint) MoneyValue {
	return MoneyValue(value)
}

var FloatRegexp = regexp.MustCompile(`^-?(\d+\.?\d*|\.\d+)([eE][-+]?\d+)?$`)

// NewMoneyValueFromString creates a new MoneyValue from a string
func NewMoneyValueFromString(value string) (MoneyValue, error) {
	// check if a decimal point is present
	if len(value) == 0 {
		return 0, fmt.Errorf("invalid value")
	}

	// check if a decimal point is present
	if FloatRegexp.MatchString(value) {
		// is float
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return 0, err
		}

		return NewDollarValue(f), nil
	}

	// is int
	f, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}

	return NewCentValue(uint(f)), nil
}

// Float64 returns the value as a float64
func (c MoneyValue) Float64() float64 {
	return float64(c) / 100.0
}

// String returns the value as a string
func (c MoneyValue) String() string {
	return fmt.Sprintf("%.2f", c.Float64())
}

func (c MoneyValue) CentsValue() uint {
	return uint(c)
}

func (u MoneyValue) GetGraphQLType() string {
	return "CentValue"
}
