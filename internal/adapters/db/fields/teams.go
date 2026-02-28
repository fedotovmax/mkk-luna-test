package fields

import (
	"errors"
	"fmt"
)

type TeamField string

const (
	TeamFieldID   TeamField = "id"
	TeamFieldName TeamField = "name"
)

func (ue TeamField) String() string {
	return string(ue)
}

var ErrInvalidTeamField = errors.New("the passed field does not belong to the team entity")

func IsTeamEntityField(f TeamField) error {

	const op = "db.fields.IsTeamEntityField"

	switch f {
	case TeamFieldName, TeamFieldID:
		return nil
	}

	return fmt.Errorf("%s: %w", op, ErrInvalidTeamField)
}
