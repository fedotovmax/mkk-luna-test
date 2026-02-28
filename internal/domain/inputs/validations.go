package inputs

import (
	"regexp"

	"github.com/fedotovmax/mkk-luna-test/internal/validation"
)

var UpperLettersRegexp = regexp.MustCompile(`[A-Z]`)
var LowerLettersRegexp = regexp.MustCompile(`[a-z]`)
var DigitRegexp = regexp.MustCompile(`\d`)
var SpecialRegexp = regexp.MustCompile(`[!@#$%^&*()_\-+=\[\]{}|\\;:'",.<>/?]`)
var NameRegexp = regexp.MustCompile(`^\p{L}+(?:[ '’]\p{L}+)*$`)

func validatePassword(password string) (string, error) {

	err := validation.MinLength(password, 8)

	if err != nil {
		return "Длина пароля должна быть минимум 8 символов", err
	}

	err = validation.Regex(password, UpperLettersRegexp)
	if err != nil {
		return "Пароль должен содержать минимум 1 символ в верхнем регистре", err
	}

	err = validation.Regex(password, LowerLettersRegexp)

	if err != nil {

		return "Пароль должен содержать минимум 1 символ в нижнем регистре", err
	}

	err = validation.Regex(password, DigitRegexp)

	if err != nil {
		return "Пароль должен содержать минимум 1 цифру", err
	}

	err = validation.Regex(password, SpecialRegexp)

	if err != nil {
		return "Пароль должен содержать минимум 1 специальный символ", err
	}

	return "", nil
}
