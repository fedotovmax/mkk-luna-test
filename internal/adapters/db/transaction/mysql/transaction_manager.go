package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/fedotovmax/mkk-luna-test/internal/adapters/db"
	"github.com/fedotovmax/mkk-luna-test/internal/adapters/db/transaction"
)

type ctxKey struct{}

type txOwner interface {
	db.StdSQLDriverQueryExecutor
	db.StdSQLDriverTx
}

type manager struct {
	pool txOwner
	log  *slog.Logger
}

type mysqlTranscation struct {
	*sql.Tx
}

func Init(conn txOwner, l ...*slog.Logger) (Manager, error) {
	if conn == nil {
		return nil, transaction.ErrConnRequired
	}
	return &manager{
		pool: conn,
		log:  l[0],
	}, nil
}

type Extractor interface {
	ExtractTx(ctx context.Context) db.StdSQLDriverQueryExecutor
}

type Manager interface {
	Wrap(ctx context.Context, fn func(context.Context) error) error
	GetExtractor() Extractor
}

func (m *manager) GetExtractor() Extractor {
	return m
}

func (m *manager) Wrap(ctx context.Context, fn func(context.Context) error) error {
	tx, err := m.pool.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("pool.Begin: cannot start transaction: %w", err)
	}

	defer func() {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			if m.log != nil {
				m.log.Error("error when try to rollback transaction", slog.String("error", rollbackErr.Error()))
			}
		} else if rollbackErr == nil {
			if m.log != nil {
				m.log.Info("transaction successfully rollbacked")
			}
		}
	}()

	ctx = context.WithValue(ctx, ctxKey{}, &mysqlTranscation{tx})

	err = fn(ctx)

	if err != nil {
		return fmt.Errorf("error when execute transaction fn: %w", err)
	}

	err = tx.Commit()

	if err != nil {
		return fmt.Errorf("error when commit: %w", err)
	}

	return nil
}

func (m *manager) ExtractTx(ctx context.Context) db.StdSQLDriverQueryExecutor {
	executor, ok := ctx.Value(ctxKey{}).(*mysqlTranscation)
	if !ok {
		return m.pool
	}

	return executor
}
