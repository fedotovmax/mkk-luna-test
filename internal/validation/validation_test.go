package validation

import (
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmptyString(t *testing.T) {
	assert.ErrorIs(t, EmptyString(""), ErrEmptyString)
	assert.NoError(t, EmptyString("not empty"))
	assert.NoError(t, EmptyString(" "))
}

func TestRange(t *testing.T) {
	t.Run("range integers", func(t *testing.T) {
		assert.NoError(t, Range(10, 1, 20))
		assert.ErrorIs(t, Range(25, 1, 20), ErrNumberOutOfRange)
		assert.ErrorIs(t, Range(0, 1, 20), ErrNumberOutOfRange)
	})

	t.Run("range floats", func(t *testing.T) {
		assert.NoError(t, Range(15.5, 10.0, 20.0))
		assert.ErrorIs(t, Range(5.5, 10.0, 20.0), ErrNumberOutOfRange)
	})
}

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

func TestMinLength(t *testing.T) {
	t.Run("valid length", func(t *testing.T) {
		assert.NoError(t, MinLength("hello", 3))
	})

	t.Run("equal length", func(t *testing.T) {
		assert.NoError(t, MinLength("hello", 5))
	})

	t.Run("too short", func(t *testing.T) {
		assert.ErrorIs(t, MinLength("hi", 3), ErrStringTooShort)
	})

	t.Run("utf8 symbols check", func(t *testing.T) {
		// 2 руны, а не 4 байта
		assert.NoError(t, MinLength("пр", 2))
		assert.ErrorIs(t, MinLength("пр", 3), ErrStringTooShort)
	})
}

func TestMaxLength(t *testing.T) {
	t.Run("valid length", func(t *testing.T) {
		assert.NoError(t, MaxLength("hello", 10))
	})

	t.Run("equal length", func(t *testing.T) {
		assert.NoError(t, MaxLength("hello", 5))
	})

	t.Run("too long", func(t *testing.T) {
		assert.ErrorIs(t, MaxLength("hello", 3), ErrStringTooLong)
	})

	t.Run("utf8 support", func(t *testing.T) {
		assert.NoError(t, MaxLength("пр", 2))
		assert.ErrorIs(t, MaxLength("привет", 3), ErrStringTooLong)
	})
}

func TestLengthRange(t *testing.T) {
	t.Run("within range", func(t *testing.T) {
		assert.NoError(t, LengthRange("hello", 3, 10))
	})

	t.Run("equal boundaries", func(t *testing.T) {
		assert.NoError(t, LengthRange("hello", 5, 5))
	})

	t.Run("too short", func(t *testing.T) {
		assert.ErrorIs(t, LengthRange("hi", 3, 5), ErrStringLengthOut)
	})

	t.Run("too long", func(t *testing.T) {
		assert.ErrorIs(t, LengthRange("hello world", 3, 5), ErrStringLengthOut)
	})
}

func TestRegex(t *testing.T) {
	re := regexp.MustCompile(`^[a-z]+$`)

	t.Run("match", func(t *testing.T) {
		assert.NoError(t, Regex("hello", re))
	})

	t.Run("no match", func(t *testing.T) {
		assert.ErrorIs(t, Regex("hello123", re), ErrStringRegexMatch)
	})

	t.Run("empty string no match", func(t *testing.T) {
		assert.ErrorIs(t, Regex("", re), ErrStringRegexMatch)
	})
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
