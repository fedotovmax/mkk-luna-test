package validation

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestEmptyString(t *testing.T) {
	assert.ErrorIs(t, EmptyString(""), ErrEmptyString)
	assert.NoError(t, EmptyString("not empty"))
	assert.NoError(t, EmptyString(" "))
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
