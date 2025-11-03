package domain

import "strings"

type ID string

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

func (id ID) String() string {
	return string(id)
}

type IDGenerator interface {
	NewID() (ID, error)
}

const (
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
