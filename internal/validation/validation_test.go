package validation

import (
	"testing"

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
