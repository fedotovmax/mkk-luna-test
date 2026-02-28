package users

import (
	"context"
	"fmt"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters/db/fields"
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
)

func (u *user) FindOne(ctx context.Context, field fields.UserField, value string) (*domain.User, error) {

	const op = "adapters.db.mysql.users.FindOne"

	err := fields.IsUserEntityField(field)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return nil, nil
}
