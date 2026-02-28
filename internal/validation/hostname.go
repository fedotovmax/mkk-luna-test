package validation

import (
	"errors"
	"strings"
)

var (
	ErrInvalidHostname = errors.New("invalid hostname")
)

func IsHostname(host string) error {
	normalized := strings.ToLower(strings.TrimSuffix(host, "."))

	if len(normalized) > 253 {
		return ErrInvalidHostname
	}

	parts := strings.Split(normalized, ".")
	for _, part := range parts {
		l := len(part)

		if l == 0 {
			return ErrInvalidHostname
		}

		if l > 63 {
			return ErrInvalidHostname
		}

		if part[0] == '-' {
			return ErrInvalidHostname
		}

		if part[l-1] == '-' {
			return ErrInvalidHostname
		}

		for _, r := range part {
			if (r < 'a' || r > 'z') &&
				(r < '0' || r > '9') &&
				r != '-' {

				return ErrInvalidHostname
			}
		}
	}

	return nil
}
