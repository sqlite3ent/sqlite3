// Package sqlite3 implements the functions, types, and interfaces for the module.
package sqlite3

import (
	"context"
	"database/sql/driver"

	"modernc.org/sqlite"
)

// sqliteConn is the interface that wraps the basic modernc.org/sqlite.conn methods.
type sqliteConn interface {
	FileControlPersistWAL(dbName string, mode int) (int, error)
	Ping(ctx context.Context) (err error)
	BeginTx(ctx context.Context, opts driver.TxOptions) (dt driver.Tx, err error)
	PrepareContext(ctx context.Context, query string) (ds driver.Stmt, err error)
	ExecContext(ctx context.Context, query string, args []driver.NamedValue) (dr driver.Result, err error)
	QueryContext(ctx context.Context, query string, args []driver.NamedValue) (dr driver.Rows, err error)
	Begin() (dt driver.Tx, err error)
	Close() (err error)
	ResetSession(ctx context.Context) error
	IsValid() bool
	Exec(query string, args []driver.Value) (dr driver.Result, err error)
	Prepare(query string) (ds driver.Stmt, err error)
	Query(query string, args []driver.Value) (dr driver.Rows, err error)
	Serialize() (v []byte, err error)
	Deserialize(buf []byte) (err error)
	NewBackup(dstUri string) (*sqlite.Backup, error)
	NewRestore(srcUri string) (*sqlite.Backup, error)
}
