package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters"
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
	"github.com/redis/go-redis/v9"
)

func (r *tasks) Get(ctx context.Context, key string) (*domain.FindTasksResponse, error) {

	const op = "adapters.cache.redis.tasks.Get"

	bytes, err := r.rdb.Get(ctx, key).Bytes()

	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("%s: %w: %v", op, adapters.ErrNotFound, err)
		}
		return nil, fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	var res domain.FindTasksResponse

	err = json.Unmarshal(bytes, &res)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &res, nil
}
