package utilities

import (
	"fmt"
	"strings"
)

type Recurrence string

func errInvalidRecurrence(val string) error {
	return fmt.Errorf("invalid recurrence.  expecting one of (daily, weekly, monthly, yearly). received: %s", val)
}

func NewRecurrence(val string) (*Recurrence, error) {
	val = strings.ToLower(val)

	var res Recurrence

	switch val {
	case "daily":
		res = Recurrence(val)
	case "weekly":
		res = Recurrence(val)
	case "monthly":
		res = Recurrence(val)
	case "yearly":
		res = Recurrence(val)
	default:
		return nil, errInvalidRecurrence(val)
	}

	return &res, nil
}

func (r Recurrence) String() string {
	return string(r)
}

func (r Recurrence) GetGraphQLType() string {
	return "Recurrence"
}
