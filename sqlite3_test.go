package sqlite3

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDriver(t *testing.T) {
	// Create a temporary file for the database
	tmpfile, err := os.CreateTemp("", "testdb-*.db")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	// Test opening a database
	db, err := sql.Open("sqlite3", tmpfile.Name())
	require.NoError(t, err)
	defer db.Close()

	// Test basic SQL operations
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS test (id INTEGER PRIMARY KEY, name TEXT)`)
	require.NoError(t, err)

	// Test insert
	result, err := db.Exec(`INSERT INTO test (name) VALUES (?)`, "test")
	require.NoError(t, err)

	// Check rows affected
	affected, err := result.RowsAffected()
	require.NoError(t, err)
	require.Equal(t, int64(1), affected)

	// Test query
	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM test`).Scan(&count)
	require.NoError(t, err)
	require.Equal(t, 1, count)
}

func TestDriverWithContext(t *testing.T) {
	// Create a temporary file for the database
	tmpfile, err := os.CreateTemp("", "testdb-ctx-*.db")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	// Test opening a database with context
	ctx := context.Background()
	db, err := sql.Open("sqlite3", tmpfile.Name())
	require.NoError(t, err)
	defer db.Close()

	// Create table with context
	_, err = db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS test_ctx (id INTEGER PRIMARY KEY, name TEXT)`)
	require.NoError(t, err)

	// Test transaction with context
	tx, err := db.BeginTx(ctx, nil)
	require.NoError(t, err)

	_, err = tx.ExecContext(ctx, `INSERT INTO test_ctx (name) VALUES (?)`, "test")
	require.NoError(t, err)

	err = tx.Commit()
	require.NoError(t, err)

	// Verify data with context
	var name string
	err = db.QueryRowContext(ctx, `SELECT name FROM test_ctx WHERE id = ?`, 1).Scan(&name)
	require.NoError(t, err)
	require.Equal(t, "test", name)
}

func TestVersion(t *testing.T) {
	// Test that the version is set and matches the expected format
	require.Regexp(t, `^v\d+\.\d+\.\d+$`, Version, "Version should be in semver format (e.g., v1.0.0)")
}

func TestSQLiteDriver_Open(t *testing.T) {
	// Create a temporary file for testing
	tmpfile, err := os.CreateTemp("", "testdb-*.db")
	require.NoError(t, err)
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	tests := []struct {
		name    string
		dsn     string
		wantErr bool
	}{
		{
			name:    "valid dsn with WAL journal mode",
			dsn:     "file:" + tmpfile.Name() + "?cache=shared&_journal=WAL&_fk=1",
			wantErr: false,
		},
		{
			name:    "invalid dsn",
			dsn:     "invalid:dsn:format",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &SQLiteDriver{}
			conn, err := d.Open(tt.dsn)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Open() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && conn == nil {
				t.Error("expected connection, got nil")
			}
			if tt.wantErr && conn != nil {
				t.Error("expected nil connection on error")
			}

			// If connection is successful, test basic operations
			if conn != nil {
				// Test Prepare
				stmt, err := conn.Prepare("SELECT 1")
				if err != nil {
					t.Errorf("Prepare failed: %v", err)
				}
				if stmt != nil {
					// Test stmt Close
					if err := stmt.Close(); err != nil {
						t.Errorf("stmt.Close() error = %v", err)
					}
				}

				// Test Begin
				tx, err := conn.Begin()
				if err != nil {
					t.Errorf("Begin failed: %v", err)
				}
				if tx != nil {
					// Test tx Rollback
					if err := tx.Rollback(); err != nil {
						t.Errorf("Rollback failed: %v", err)
					}
				}

				// Test Close
				if err := conn.Close(); err != nil {
					t.Errorf("Close failed: %v", err)
				}
			}
		})
	}
}
