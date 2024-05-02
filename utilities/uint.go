package utilities

import "fmt"

type UInt uint

func NewUInt(val uint) UInt {
	return UInt(val)
}

func (u UInt) GetGraphQLType() string {
	return "UInt"
}

func (u UInt) String() string {
	return fmt.Sprintf("%d", u)
}

func (u UInt) Uint() uint {
	return uint(u)
}

func NewUIntPointer(val *uint) *UInt {
	if val == nil {
		return nil
	}

	result := NewUInt(*val)
	return &result
}
