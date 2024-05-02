package utilities

import (
	"fmt"
	"regexp"

	"github.com/google/uuid"
)

type AltUUID string

func errInvalidAltUUID(val string) error {
	return fmt.Errorf("invalid alt id or uuid.  expecting format like (12345678-1234-1234-1234-123456789abc) or (my_alt_id). received: %s", val)
}

func NewAltUUID(altIdOrUuid string) (AltUUID, error) {
	result, err := NewAltUUIDPointer(&altIdOrUuid)
	if result == nil {
		return "", err
	}
	return *result, err
}
func NewAltUUIDPointer(altIdOrUuid *string) (*AltUUID, error) {
	if altIdOrUuid == nil {
		return nil, nil
	}

	// check if is uuid
	if _, err := uuid.Parse(*altIdOrUuid); err != nil {
		// is not uuid, check if lower, snake case
		if !isLowerSnakeCase(*altIdOrUuid) {
			return nil, errInvalidAltUUID(*altIdOrUuid)
		}
	}

	parsed := AltUUID(*altIdOrUuid)
	return &parsed, nil
}

func (u AltUUID) String() string {
	return string(u)
}

func (u AltUUID) GetGraphQLType() string {
	return "AltUuid"
}

func isLowerSnakeCase(s string) bool {
	match, _ := regexp.MatchString(`^[a-z]+(_[a-z]+)*$`, s)
	return match
}
