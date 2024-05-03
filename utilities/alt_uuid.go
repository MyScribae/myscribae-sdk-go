package utilities

import (
	"fmt"
	"regexp"

	"github.com/google/uuid"
)

type AltUuid string

func errInvalidAltUuid(val string) error {
	return fmt.Errorf("invalid alt id or uuid.  expecting format like (12345678-1234-1234-1234-123456789abc) or (my_alt_id). received: %s", val)
}

func NewAltUuid(altIdOrUuid string) (AltUuid, error) {
	result, err := NewAltUuidPointer(&altIdOrUuid)
	if result == nil {
		return "", err
	}
	return *result, err
}
func NewAltUuidPointer(altIdOrUuid *string) (*AltUuid, error) {
	if altIdOrUuid == nil {
		return nil, nil
	}

	// check if is uuid
	if _, err := uuid.Parse(*altIdOrUuid); err != nil {
		// is not uuid, check if lower, snake case
		if !isLowerSnakeCase(*altIdOrUuid) {
			return nil, errInvalidAltUuid(*altIdOrUuid)
		}
	}

	parsed := AltUuid(*altIdOrUuid)
	return &parsed, nil
}

func (u AltUuid) String() string {
	return string(u)
}

func (u *AltUuid) GetGraphQLType() string {
	return "AltUuid"
}

func isLowerSnakeCase(s string) bool {
	match, _ := regexp.MatchString(`^[a-z]+(_[a-z]+)*$`, s)
	return match
}
