package validation

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsHostname(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{"simple domain", "example.com", nil},
		{"subdomain", "mail.example.com", nil},
		{"uppercase normalized", "EXAMPLE.COM", nil},
		{"trailing dot", "example.com.", nil},
		{"numbers allowed", "api2.example3.com", nil},
		{"hyphen inside", "my-site.example.com", nil},

		{"empty string", "", ErrInvalidHostname},
		{"label too long", strings.Repeat("a", 64) + ".com", ErrInvalidHostname},
		{"domain too long", strings.Repeat("a.", 127) + "a", ErrInvalidHostname},
		{"label starts with hyphen", "-example.com", ErrInvalidHostname},
		{"label ends with hyphen", "example-.com", ErrInvalidHostname},
		{"invalid symbol", "exa$mple.com", ErrInvalidHostname},
		{"double dot", "example..com", ErrInvalidHostname},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := IsHostname(tt.input)

			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}
