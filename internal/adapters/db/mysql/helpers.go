package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/fedotovmax/mkk-luna-test/pkg/logger"
)

func constructorCloseConnection(db *sql.DB, log *slog.Logger) {
	if db != nil {
		err := db.Close()
		if err != nil {
			log.Error("failed when close mysql database connection", logger.Err(err))
		}
	}
}

func connectWithRetries(

	ctx context.Context,

	log *slog.Logger,

	dsn string,

	maxRetries uint8,

	retryWait time.Duration,

) (*sql.DB, error) {

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	var lastPingError error

	for i := uint8(1); i <= maxRetries; i++ {
		if ctx.Err() != nil {
			constructorCloseConnection(db, log)
			return nil, ctx.Err()
		}
		lastPingError = db.PingContext(ctx)

		if lastPingError == nil {
			return db, nil
		}

		log.Warn("mysql ping failed", slog.Int("attempt", int(i)), logger.Err(lastPingError))

		select {
		case <-time.After(retryWait):
		case <-ctx.Done():
			constructorCloseConnection(db, log)
			return nil, ctx.Err()
		}
	}
	constructorCloseConnection(db, log)
	return nil, fmt.Errorf("connection to mysql failed after %d attempts: %w", maxRetries, lastPingError)
}
