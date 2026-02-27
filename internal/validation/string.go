package validation

import (
	"errors"
	"regexp"
	"unicode/utf8"
)

var ErrEmptyString = errors.New("string is empty")

var ErrStringTooShort = errors.New("string is shorter than minimum length")

var ErrStringTooLong = errors.New("string is longer than maximum length")

var ErrStringLengthOut = errors.New("string length out of allowed range")

var ErrStringRegexMatch = errors.New("string does not match required pattern")

func EmptyString(s string) error {
	if s == "" {
		return ErrEmptyString
	}
	return nil
}

func MinLength(s string, min int) error {
	l := utf8.RuneCountInString(s)
	if l < min {
		return ErrStringTooShort
	}
	return nil
}

func MaxLength(s string, max int) error {
	l := utf8.RuneCountInString(s)
	if l > max {
		return ErrStringTooLong
	}
	return nil
}

func LengthRange(s string, min, max int) error {
	l := utf8.RuneCountInString(s)
	if l < min || l > max {
		return ErrStringLengthOut
	}
	return nil
}

func Regex(s string, re *regexp.Regexp) error {
	if !re.MatchString(s) {
		return ErrStringRegexMatch
	}
	return nil
}
