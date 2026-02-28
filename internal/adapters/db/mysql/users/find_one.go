package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters"
	"github.com/fedotovmax/mkk-luna-test/internal/adapters/db/fields"
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
)

func (u *user) FindOne(ctx context.Context, field fields.UserField, value string) (*domain.User, error) {

	const op = "adapters.db.mysql.users.find_one"

	err := fields.IsUserEntityField(field)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	tx := u.txExtractor.ExtractTx(ctx)

	row := tx.QueryRowContext(ctx, findByQuery(field), value)

	user := &domain.User{}

	err = row.Scan(
		&user.ID,
		&user.UserName,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %s: %w: %v", op, field, adapters.ErrNotFound, err)
		}
		return nil, fmt.Errorf("%s: %s: %w: %v", op, field, adapters.ErrInternal, err)
	}

	return user, nil
}
