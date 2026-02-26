package validation

import "errors"

var ErrEmptyString = errors.New("string is empty")

func EmptyString(s string) error {
	if s == "" {
		return ErrEmptyString
	}
	return nil
}
