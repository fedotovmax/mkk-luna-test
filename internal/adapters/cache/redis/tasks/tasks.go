package tasks

import (
	"log/slog"

	"github.com/redis/go-redis/v9"
)

type tasks struct {
	log *slog.Logger
	rdb *redis.Client
}

func New(log *slog.Logger, rdb *redis.Client) *tasks {
	return &tasks{
		log: log,
		rdb: rdb,
	}
}
