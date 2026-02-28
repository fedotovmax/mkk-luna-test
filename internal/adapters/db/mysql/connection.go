package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters/db"
	"github.com/fedotovmax/mkk-luna-test/internal/config"
	mysqlDriver "github.com/go-sql-driver/mysql"
)

type pool struct {
	log *slog.Logger
	*sql.DB
}

var (
	mysqlOnce sync.Once
	mysqlPool *sql.DB
	initErr   error
)

func New(

	ctx context.Context,

	log *slog.Logger,

	cfg *config.Database,
) (db.StdSQLDriver, error) {

	const op = "adapters.db.mysql.New"

	l := log.With(slog.String("op", op))

	mysqlOnce.Do(func() {

		_, err := mysqlDriver.ParseDSN(cfg.DSN)

		if err != nil {
			initErr = fmt.Errorf("%s: %w: %v", op, ErrInvalidDSNFormat, err)
			return
		}

		mysqlPool, initErr = connectWithRetries(ctx, l, cfg.DSN, cfg.MaxRetries, cfg.RetryWait)
		if initErr == nil && mysqlPool != nil {
			mysqlPool.SetMaxOpenConns(15)
			mysqlPool.SetMaxIdleConns(5)
			mysqlPool.SetConnMaxLifetime(30 * time.Minute)
			mysqlPool.SetConnMaxIdleTime(5 * time.Minute)
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
