package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters"
	"github.com/fedotovmax/mkk-luna-test/internal/domain"
)

func (t *tasks) Set(ctx context.Context, key string, data *domain.FindTasksResponse) error {

	const op = "adapters.cache.redis.tasks.Set"

	bytes, err := json.Marshal(data)

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	err = t.rdb.Set(ctx, key, bytes, time.Minute*5).Err()

	if err != nil {
		return fmt.Errorf("%s: %w: %v", op, adapters.ErrInternal, err)
	}

	return nil
}
