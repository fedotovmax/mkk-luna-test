package fields

import (
	"errors"
	"fmt"
)

type UserField string

const (
	UserFieldUsername UserField = "username"
	UserFieldEmail    UserField = "email"
	UserFieldID       UserField = "id"
)

func (ue UserField) String() string {
	return string(ue)
}

var ErrInvalidUserField = errors.New("the passed field does not belong to the user entity")

func IsUserEntityField(f UserField) error {

	const op = "db.fields.IsUserEntityField"

	switch f {
	case UserFieldUsername, UserFieldEmail, UserFieldID:
		return nil
	}

	return fmt.Errorf("%s: %w", op, ErrInvalidUserField)
}
