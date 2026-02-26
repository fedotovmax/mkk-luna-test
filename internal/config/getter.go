package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

var ErrVariableParse = errors.New("error when parse variable from env")

var ErrVariableNotProvided = errors.New("env variable not provided")

var ErrVariableIsEmpty = errors.New("env variable is empty")

var ErrUnsupportedType = errors.New("unsupported variable type")

func getEnv(name string) (string, error) {
	const op = "config.getEnv"
	value, exists := os.LookupEnv(name)

	if !exists {
		return "", fmt.Errorf("%s: variable key: %s: %w", op, name, ErrVariableNotProvided)
	}

	if value == "" {
		return "", fmt.Errorf("%s: variable key: %s: %w", op, name, ErrVariableIsEmpty)
	}

	return value, nil
}

func getEnvAs[T any](name string) (T, error) {
	const op = "config.getEnvAs"

	var zero T

	valueStr, err := getEnv(name)
	if err != nil {
		return zero, err
	}

	value, err := parseValue[T](valueStr)

	if err != nil {
		return zero, fmt.Errorf("%s: variable key: %s: %w", op, name, err)
	}

	return value, nil
}

func parseValue[T any](v string) (T, error) {
	var zero T
	var t T

	switch any(t).(type) {

	case int:
		i, err := strconv.Atoi(v)
		if err != nil {
			return zero, fmt.Errorf("%w: %v", ErrVariableParse, err)
		}
		return any(i).(T), nil

	case uint8:
		u, err := strconv.ParseUint(v, 10, 8)
		if err != nil {
			return zero, fmt.Errorf("%w: %v", ErrVariableParse, err)
		}
		return any(uint8(u)).(T), nil

	case uint16:
		u, err := strconv.ParseUint(v, 10, 16)
		if err != nil {
			return zero, fmt.Errorf("%w: %v", ErrVariableParse, err)
		}
		return any(uint16(u)).(T), nil

	case uint32:
		u, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return zero, fmt.Errorf("%w: %v", ErrVariableParse, err)
		}
		return any(uint32(u)).(T), nil

	case uint64:
		u, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return zero, fmt.Errorf("%w: %v", ErrVariableParse, err)
		}
		return any(uint64(u)).(T), nil

	case float32:
		f, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return zero, fmt.Errorf("%w: %v", ErrVariableParse, err)
		}
		return any(float32(f)).(T), nil

	case float64:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return zero, fmt.Errorf("%w: %v", ErrVariableParse, err)
		}
		return any(f).(T), nil

	case bool:
		b, err := strconv.ParseBool(v)
		if err != nil {
			return zero, fmt.Errorf("%w: %v", ErrVariableParse, err)
		}
		return any(b).(T), nil

	case string:
		return any(v).(T), nil

	case time.Duration:
		dur, err := time.ParseDuration(v)
		if err != nil {
			return zero, fmt.Errorf("%w: %v", ErrVariableParse, err)
		}

		return any(dur).(T), nil
	}

	return zero, ErrUnsupportedType
}
