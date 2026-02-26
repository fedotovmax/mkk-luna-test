package validation

import "errors"

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
	ErrNumberOutOfRange = errors.New("number is outside allowed range")
)

func Range[T Numbers](value, min, max T) error {
	if value < min || value > max {
		return ErrNumberOutOfRange
	}
	return nil
}
