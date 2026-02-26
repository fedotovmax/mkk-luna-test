package db

import (
	"context"
	"database/sql"
)

type StdSQLDriverTx interface {
	Begin() (*sql.Tx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type StdSQLDriver interface {
	StdSQLDriverTx
	StdSQLDriverQueryExecutor
	SQLDatabaseLifecycle
}

type StdSQLDriverQueryExecutor interface {
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type SQLDatabaseLifecycle interface {
	Stop(ctx context.Context) error
}
