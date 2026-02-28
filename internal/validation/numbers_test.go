package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestMin(t *testing.T) {
	t.Run("int valid", func(t *testing.T) {
		assert.NoError(t, Min(10, 5))
	})

	t.Run("int equal", func(t *testing.T) {
		assert.NoError(t, Min(5, 5))
	})

	t.Run("int too small", func(t *testing.T) {
		assert.ErrorIs(t, Min(3, 5), ErrNumberTooSmall)
	})

	t.Run("float valid", func(t *testing.T) {
		assert.NoError(t, Min(10.5, 5.2))
	})

	t.Run("float too small", func(t *testing.T) {
		assert.ErrorIs(t, Min(1.1, 2.2), ErrNumberTooSmall)
	})
}

func TestMax(t *testing.T) {
	t.Run("int valid", func(t *testing.T) {
		assert.NoError(t, Max(5, 10))
	})

	t.Run("int equal", func(t *testing.T) {
		assert.NoError(t, Max(10, 10))
	})

	t.Run("int too large", func(t *testing.T) {
		assert.ErrorIs(t, Max(15, 10), ErrNumberTooLarge)
	})

	t.Run("float valid", func(t *testing.T) {
		assert.NoError(t, Max(5.5, 10))
	})

	t.Run("float too large", func(t *testing.T) {
		assert.ErrorIs(t, Max(15.5, 10), ErrNumberTooLarge)
	})
}

func TestEqual(t *testing.T) {
	t.Run("int equal", func(t *testing.T) {
		assert.NoError(t, Equal(5, 5))
	})

	t.Run("int not equal", func(t *testing.T) {
		assert.ErrorIs(t, Equal(5, 10), ErrNumberNotEqual)
	})

	t.Run("float equal", func(t *testing.T) {
		assert.NoError(t, Equal(5.5, 5.5))
	})

	t.Run("float not equal", func(t *testing.T) {
		assert.ErrorIs(t, Equal(5.5, 5.6), ErrNumberNotEqual)
	})
}

func TestIsPositive(t *testing.T) {
	t.Run("positive int", func(t *testing.T) {
		assert.NoError(t, IsPositive(10))
	})

	t.Run("zero int", func(t *testing.T) {
		assert.ErrorIs(t, IsPositive(0), ErrNumberNotPositive)
	})

	t.Run("negative int", func(t *testing.T) {
		assert.ErrorIs(t, IsPositive(-5), ErrNumberNotPositive)
	})

	t.Run("positive float", func(t *testing.T) {
		assert.NoError(t, IsPositive(3.14))
	})

	t.Run("negative float", func(t *testing.T) {
		assert.ErrorIs(t, IsPositive(-1.1), ErrNumberNotPositive)
	})
}

func TestIsNegative(t *testing.T) {
	t.Run("negative int", func(t *testing.T) {
		assert.NoError(t, IsNegative(-10))
	})

	t.Run("zero int", func(t *testing.T) {
		assert.ErrorIs(t, IsNegative(0), ErrNumberNotNegative)
	})

	t.Run("positive int", func(t *testing.T) {
		assert.ErrorIs(t, IsNegative(5), ErrNumberNotNegative)
	})

	t.Run("negative float", func(t *testing.T) {
		assert.NoError(t, IsNegative(-3.14))
	})

	t.Run("positive float", func(t *testing.T) {
		assert.ErrorIs(t, IsNegative(1.1), ErrNumberNotNegative)
	})
}
