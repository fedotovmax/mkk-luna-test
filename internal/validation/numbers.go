package validation

import (
	"errors"
)

type Float interface {
	~float32 | ~float64
}

type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Integer interface {
	Signed | Unsigned
}

type Numbers interface {
	Integer | Float
}

var (
	ErrNumberOutOfRange  = errors.New("number is outside allowed range")
	ErrNumberTooSmall    = errors.New("number is less than minimum allowed")
	ErrNumberTooLarge    = errors.New("number exceeds maximum allowed")
	ErrNumberNotEqual    = errors.New("number is not equal to required value")
	ErrNumberNotPositive = errors.New("number must be positive")
	ErrNumberNotNegative = errors.New("number must be negative")
)

func Range[T Numbers](value, min, max T) error {
	if value < min || value > max {
		return ErrNumberOutOfRange
	}
	return nil
}

func Min[T Numbers](value, min T) error {
	if value < min {
		return ErrNumberTooSmall
	}
	return nil
}

func Max[T Numbers](value, max T) error {
	if value > max {
		return ErrNumberTooLarge
	}
	return nil
}

func Equal[T Numbers](value, expected T) error {
	if value != expected {
		return ErrNumberNotEqual
	}
	return nil
}

func IsPositive[T Numbers](value T) error {
	if value <= 0 {
		return ErrNumberNotPositive
	}
	return nil
}

func IsNegative[T Numbers](value T) error {
	if value >= 0 {
		return ErrNumberNotNegative
	}
	return nil
}
