package redis

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/fedotovmax/mkk-luna-test/internal/config"
	goredis "github.com/redis/go-redis/v9"
)

type RedisDb struct {
	redisClient *goredis.Client
	log         *slog.Logger
}

func New(ctx context.Context, cfg *config.Redis, log *slog.Logger) (*RedisDb, error) {

	const op = "adapters.cache.redis.New"

	redisClient := goredis.NewClient(&goredis.Options{
		Addr:            cfg.Addr,
		Password:        cfg.Password,
		DB:              cfg.DB,
		MaxRetries:      int(cfg.MaxRetries),
		MinRetryBackoff: cfg.RetryWait,
		MaxRetryBackoff: cfg.RetryWait,
		PoolSize:        20,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour * 1,
		ConnMaxIdleTime: time.Minute * 10,
	})

	_, err := redisClient.Ping(ctx).Result()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &RedisDb{
		redisClient: redisClient,
		log:         log,
	}, nil

}

func (r *RedisDb) Stop(ctx context.Context) error {
	op := "adapters.cache.redis.Stop"

	done := make(chan error, 1)

	go func() {
		err := r.redisClient.Close()
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("%s: %w: %v", op, ErrCloseTimeout, ctx.Err())
	}
}
