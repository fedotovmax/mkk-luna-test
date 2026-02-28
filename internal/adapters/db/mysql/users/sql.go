package users

import (
	"fmt"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters/db/fields"
)

const create = "insert into users (id, username, email, password_hash) values (?, ?, ?, ?);"

func findByQuery(field fields.UserField) string {
	return fmt.Sprintf("select id, username, email, password_hash, created_at, updated_at from users where %s = ?;", field)
}
