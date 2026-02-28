package validation

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmail(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{

		{"valid simple", "user@example.com", nil},
		{"valid uppercase domain", "user@EXAMPLE.COM", nil},
		{"valid with subdomain", "user@mail.example.com", nil},
		{"valid with display name", "John Doe <john@example.com>", nil},

		{"missing at", "userexample.com", ErrEmailInvalidFormat},
		{"missing domain", "user@", ErrEmailInvalidFormat},
		{"empty string", "", ErrEmailInvalidFormat},

		{
			"too long email",
			strings.Repeat("a", 245) + "@example.com",
			ErrEmailTooLong,
		},

		{
			"local part too long",
			strings.Repeat("a", 65) + "@example.com",
			ErrEmailLocalPartTooLong,
		},

		{
			"invalid hostname",
			"user@-example.com",
			ErrInvalidHostname,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := IsEmail(tt.input)

			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}
