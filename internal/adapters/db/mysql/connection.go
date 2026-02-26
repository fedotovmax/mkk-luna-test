package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters/db"
	"github.com/fedotovmax/mkk-luna-test/pkg/logger"
	mysqlDriver "github.com/go-sql-driver/mysql"
)

var ErrWantToCallMethodsAfterInitPool = errors.New("you will be able to call mysql pool methods only after the connection has been created and initialized")

var ErrInvalidDSNFormat = errors.New("invalid mysql dsn format")

type pool struct {
	log *slog.Logger
	*sql.DB
}

var (
	mysqlOnce sync.Once
	mysqlPool *sql.DB
	initErr   error
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

	for i := 1; i <= 5; i++ {
		if ctx.Err() != nil {
			constructorCloseConnection(db, log)
			return nil, ctx.Err()
		}
		lastPingError = db.PingContext(ctx)

		if lastPingError == nil {
			return db, nil
		}

		log.Warn("mysql ping failed", slog.Int("attempt", i), logger.Err(lastPingError))

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

func New(

	ctx context.Context,

	log *slog.Logger,

	dsn string,

	maxRetries uint8,

	retryWait time.Duration,

) (db.StdSQLDriver, error) {

	const op = "adapters.db.mysql.New"

	l := log.With(slog.String("op", op))

	mysqlOnce.Do(func() {

		_, err := mysqlDriver.ParseDSN(dsn)

		if err != nil {
			initErr = fmt.Errorf("%s: %w: %v", op, ErrInvalidDSNFormat, err)
			return
		}

		mysqlPool, initErr = connectWithRetries(ctx, l, dsn, maxRetries, retryWait)
		if initErr == nil && mysqlPool != nil {
			mysqlPool.SetMaxOpenConns(10)
			mysqlPool.SetMaxIdleConns(5)
			mysqlPool.SetConnMaxLifetime(5 * time.Minute)
			mysqlPool.SetConnMaxIdleTime(1 * time.Minute)
		}
	})

	if initErr != nil {
		return nil, fmt.Errorf("%s: %w", op, initErr)
	}

	if mysqlPool == nil {
		return nil, fmt.Errorf("%s: mysql pool is nil after connection", op)
	}

	return &pool{log: log, DB: mysqlPool}, nil
}

func (p *pool) Stop(ctx context.Context) error {

	op := "adapters.db.mysql.Stop"

	if p == nil {
		return fmt.Errorf("%s: %w", op, ErrWantToCallMethodsAfterInitPool)
	}

	if p.DB == nil {
		return fmt.Errorf("%s: %w", op, ErrWantToCallMethodsAfterInitPool)
	}

	done := make(chan error, 1)

	go func() {
		done <- p.DB.Close()
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("%s: %w", op, ctx.Err())
	}
}
