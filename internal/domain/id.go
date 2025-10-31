package domain

import "strings"

// ID is a strongly-typed identifier for domain entities.
type ID string

// ParseID validates ULID format and returns typed identifier.
func ParseID(value string) (ID, error) {
	if value == "" {
		return "", ErrInvalidID
	}

	value = strings.ToUpper(value)

	if len(value) != 26 {
		return "", ErrInvalidID
	}

	for _, ch := range value {
		if !isValidULIDChar(ch) {
			return "", ErrInvalidID
		}
	}

	if value == zeroULID {
		return "", ErrInvalidID
	}

	return ID(value), nil
}

// String renders identifier as plain string.
func (id ID) String() string {
	return string(id)
}

// IDGenerator produces unique identifiers compatible with domain expectations.
type IDGenerator interface {
	NewID() (ID, error)
}

const (
	// ULIDAlphabet defines allowed characters for ULID identifiers.
	ULIDAlphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"
	zeroULID     = "00000000000000000000000000"
)

func isValidULIDChar(ch rune) bool {
	for _, allowed := range ULIDAlphabet {
		if ch == allowed {
			return true
		}
	}
	return false
}
