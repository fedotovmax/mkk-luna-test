package validation

import (
	"errors"
	"net/mail"
	"strings"
)

var (
	ErrEmailInvalidFormat    = errors.New("invalid email format")
	ErrEmailTooLong          = errors.New("email exceeds maximum length")
	ErrEmailLocalPartTooLong = errors.New("email local-part exceeds maximum length")
)

func IsEmail(addr string) error {
	a, err := mail.ParseAddress(addr)

	if err != nil {
		return ErrEmailInvalidFormat
	}

	addr = a.Address

	if len(addr) > 254 {
		return ErrEmailTooLong
	}

	parts := strings.SplitN(addr, "@", 2)

	if len(parts[0]) > 64 {
		return ErrEmailLocalPartTooLong
	}

	if err = IsHostname(parts[1]); err != nil {
		return err
	}

	return nil
}
