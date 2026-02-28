package validation

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsURI(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{"valid https", "https://google.com", nil},
		{"valid ftp", "ftp://files.com", nil},
		{"relative path", "/index.html", ErrURINotAbsolute},
		{"empty string", "", ErrURINotAbsolute},
		{"invalid characters", "http://a b.com", assert.AnError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uri, err := IsURI(tt.input)

			switch tt.wantErr {
			case nil:
				require.NoError(t, err)
				assert.NotNil(t, uri)
				assert.True(t, uri.IsAbs())
			case assert.AnError:
				assert.Error(t, err)
			default:
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func TestIsUUID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{
			name:    "valid uuid",
			input:   "550e8400-e29b-41d4-a716-446655440000",
			wantErr: nil,
		},
		{
			name:    "invalid format",
			input:   "not-a-uuid",
			wantErr: ErrInvalidUUIDFormat,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: ErrInvalidUUIDFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uid, err := IsUUID(tt.input)

			if tt.wantErr == nil {
				require.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, uid)
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, uuid.Nil, uid)
			}
		})
	}
}

func TestIsFilePath(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{"valid relative path", "folder/file.txt", nil},
		{"valid absolute path", "/var/log/app.log", nil},
		{"valid with spaces trimmed", "  file.txt  ", nil},

		{"empty string", "", ErrEmptyPath},
		{"only spaces", "   ", ErrEmptyPath},

		{"invalid characters 1", "file?.txt", ErrInvalidPath},
		{"invalid characters 2", "file<name>.txt", ErrInvalidPath},

		{"dot", ".", ErrInvalidPath},
		{"double dot", "..", ErrInvalidPath},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := IsFilePath(tt.input)

			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}
