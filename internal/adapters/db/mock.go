package db

import (
	"context"
	"database/sql"
)

type MockStdSQLDriver struct {
	Called   bool
	QueryStr string
	Args     []any
}

func (m *MockStdSQLDriver) Exec(query string, args ...any) (sql.Result, error) {
	m.Called = true
	m.QueryStr = query
	m.Args = args
	return sqlMockExecResult{}, nil
}

func (m *MockStdSQLDriver) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	m.Called = true
	m.QueryStr = query
	m.Args = args
	return sqlMockExecResult{}, nil
}

func (m *MockStdSQLDriver) Prepare(query string) (*sql.Stmt, error) {
	m.Called = true
	m.QueryStr = query
	return nil, nil
}

func (m *MockStdSQLDriver) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	m.Called = true
	m.QueryStr = query
	return nil, nil
}

func (m *MockStdSQLDriver) Query(query string, args ...any) (*sql.Rows, error) {
	m.Called = true
	m.QueryStr = query
	m.Args = args
	return nil, nil
}

func (m *MockStdSQLDriver) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	m.Called = true
	m.QueryStr = query
	m.Args = args
	return nil, nil
}

func (m *MockStdSQLDriver) QueryRow(query string, args ...any) *sql.Row {
	m.Called = true
	m.QueryStr = query
	m.Args = args
	return &sql.Row{}
}

func (m *MockStdSQLDriver) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	m.Called = true
	m.QueryStr = query
	m.Args = args
	return &sql.Row{}
}

type MockExtractor struct {
	Exec *MockStdSQLDriver
}

func (m *MockExtractor) ExtractTx(ctx context.Context) StdSQLDriverQueryExecutor {
	return m.Exec
}

type sqlMockExecResult struct{}

func (s sqlMockExecResult) LastInsertId() (int64, error) { return 1, nil }
func (s sqlMockExecResult) RowsAffected() (int64, error) { return 1, nil }
