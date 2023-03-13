package sqlite3

import (
	"context"
	"database/sql/driver"
)

type SQLiteConn interface {
	Ping(ctx context.Context) error
	BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error)
	PrepareContext(ctx context.Context, query string) (driver.Stmt, error)
	ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error)
	QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error)
	Begin() (driver.Tx, error)
	Close() error
	Exec(query string, args []driver.Value) (driver.Result, error)
	Prepare(query string) (driver.Stmt, error)
	Query(query string, args []driver.Value) (driver.Rows, error)
}
