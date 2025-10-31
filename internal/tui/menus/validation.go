package menus

import (
	"errors"
	"strings"
)

// ApplyValidation runs provided validator, updates field error text and reports whether an error occurred.
func ApplyValidation(screen *Screen, fieldKey, value string, validate func(string) error) bool {
	if screen == nil || validate == nil {
		return false
	}

	if err := validate(value); err != nil {
		screen.SetFieldError(fieldKey, err.Error())
		return true
	}

	screen.SetFieldError(fieldKey, "")
	return false
}

// ClearFields resets value and error state for provided keys.
func ClearFields(screen *Screen, keys ...string) {
	if screen == nil {
		return
	}

	for _, key := range keys {
		screen.SetFieldError(key, "")
		screen.SetValue(key, "")
	}
}

// ValidateNonEmpty ensures that value is not blank (after trimming).
func ValidateNonEmpty(value, message string) error {
	if strings.TrimSpace(value) == "" {
		if message == "" {
			message = "значение не может быть пустым"
		}
		return errors.New(message)
	}
	return nil
}
