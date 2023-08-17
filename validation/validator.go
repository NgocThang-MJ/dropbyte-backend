package validation

import (
	"fmt"
	"net/mail"
	"regexp"
)

var isValidFullName = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString

func ValidateLen(value string, minLen int, maxLen int) error {
	n := len(value)
	if n < minLen || n > maxLen {
		return fmt.Errorf("Must contain from %d-%d characters", minLen, maxLen)
	}
	return nil
}

func ValidateFullName(value string) error {
	if err := ValidateLen(value, 3, 100); err != nil {
		return err
	}
	if !isValidFullName(value) {
		return fmt.Errorf("Must contain only letters or spaces")
	}

	return nil
}

func ValidatePassword(value string) error {
	return ValidateLen(value, 6, 100)
}

func ValidateEmail(value string) error {
	if err := ValidateLen(value, 3, 200); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("Not a valid email address")
	}

	return nil
}
