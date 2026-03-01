package users

import (
	"context"
	"fmt"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters"
	"github.com/fedotovmax/mkk-luna-test/internal/domain/inputs"
	"github.com/google/uuid"
)

func (u *user) Create(ctx context.Context, in *inputs.CreateUser) (string, error) {

	const op = "adapters.db.mysql.users.create"

	tx := u.txExtractor.ExtractTx(ctx)

	id := uuid.New().String()

	_, err := tx.ExecContext(ctx, create, id, in.UserName, in.Email, in.Password)

	if err != nil {
		return "", fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return id, nil
}

const create = "insert into users (id, username, email, password_hash) values (?, ?, ?, ?);"
