package domain

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidatationErrors struct {
	Errors []ValidationError `json:"errors"`
}

func (e *ValidatationErrors) Error() string {

	errs := make([]string, len(e.Errors))

	for idx := range e.Errors {
		msg := fmt.Sprintf("field: %s, message: %s", e.Errors[idx].Field, e.Errors[idx].Message)
		errs = append(errs, msg)
	}

	return strings.Join(errs, ":")

}
