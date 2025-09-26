package sqlite3

import (
	"context"
	"database/sql"
	"os"
	"regexp"
	"testing"
)

func TestDriver(t *testing.T) {
	// Create a temporary file for the database
	tmpfile, err := os.CreateTemp("", "testdb-*.db")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	// Test opening a database
	db, err := sql.Open("sqlite3", tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Test basic SQL operations
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS test (id INTEGER PRIMARY KEY, name TEXT)`)
	if err != nil {
		t.Fatal(err)
	}

	// Test insert
	result, err := db.Exec(`INSERT INTO test (name) VALUES (?)`, "test")
	if err != nil {
		t.Fatal(err)
	}

	// Check rows affected
	affected, err := result.RowsAffected()
	if err != nil {
		t.Fatal(err)
	}
	if affected != int64(1) {
		t.Fatalf("expected 1 row affected, but got %d", affected)
	}

	// Test query
	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM test`).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expected count to be 1, but got %d", count)
	}
}

func TestDriverWithContext(t *testing.T) {
	// Create a temporary file for the database
	tmpfile, err := os.CreateTemp("", "testdb-ctx-*.db")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	// Test opening a database with context
	ctx := context.Background()
	db, err := sql.Open("sqlite3", tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Create table with context
	_, err = db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS test_ctx (id INTEGER PRIMARY KEY, name TEXT)`)
	if err != nil {
		t.Fatal(err)
	}

	// Test transaction with context
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO test_ctx (name) VALUES (?)`, "test")
	if err != nil {
		t.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		t.Fatal(err)
	}

	// Verify data with context
	var name string
	err = db.QueryRowContext(ctx, `SELECT name FROM test_ctx WHERE id = ?`, 1).Scan(&name)
	if err != nil {
		t.Fatal(err)
	}
	if name != "test" {
		t.Fatalf("expected name to be 'test', but got %q", name)
	}
}

func TestVersion(t *testing.T) {
	// Test that the version is set and matches the expected format
	pattern := `^v\d+\.\d+\.\d+$`
	msg := "Version should be in semver format (e.g., v1.0.0)"
	matched, err := regexp.MatchString(pattern, Version)
	if err != nil {
		t.Fatalf("failed to compile regexp: %v", err)
	}
	if !matched {
		t.Fatalf("Version %q does not match pattern %q: %s", Version, pattern, msg)
	}
}

func TestSQLiteDriver_Open(t *testing.T) {
	// Create a temporary file for testing
	tmpfile, err := os.CreateTemp("", "testdb-*.db")
	if err != nil {
		t.Fatal(err)
	}
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
