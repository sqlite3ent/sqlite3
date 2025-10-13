// Package sqlite3 implements the functions, types, and interfaces for the module.
package sqlite3

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net/url"
	"runtime"
	"strconv"
	"strings"

	"modernc.org/sqlite"
)

// This variable can be replaced with -ldflags like below:
// go build -ldflags="-X 'github.com/sqlite3ent/sqlite3.driverName=my-sqlite3'"
var driverName = "sqlite3"

func init() {
	if driverName != "" {
		sql.Register(driverName, &SQLiteDriver{})
	}
}

type SQLiteDriver struct {
	drv sqlite.Driver
}

func (d *SQLiteDriver) Open(dsn string) (driver.Conn, error) {
	var pkey string

	// COMMENT_FLAG: don't support
	// Options
	//var loc *time.Location
	//authCreate := false
	//authUser := ""
	//authPass := ""
	//authCrypt := ""
	//authSalt := ""
	//mutex := C.int(C.SQLITE_OPEN_FULLMUTEX)
	//txlock := "BEGIN"

	// PRAGMA's
	autoVacuum := -1
	busyTimeout := 5000
	caseSensitiveLike := -1
	deferForeignKeys := -1
	foreignKeys := -1
	ignoreCheckConstraints := -1
	var journalMode string
	lockingMode := "NORMAL"
	queryOnly := -1
	recursiveTriggers := -1
	secureDelete := "DEFAULT"
	synchronousMode := "NORMAL"
	writableSchema := -1
	//vfsName := ""
	var cacheSize *int64

	pos := strings.IndexRune(dsn, '?')
	if pos >= 1 {
		params, err := url.ParseQuery(dsn[pos+1:])
		if err != nil {
			return nil, err
		}

		// COMMENT_FLAG: don't support
		// Authentication
		//if _, ok := params["_auth"]; ok {
		//	authCreate = true
		//}
		//if val := params.Get("_auth_user"); val != "" {
		//	authUser = val
		//}
		//if val := params.Get("_auth_pass"); val != "" {
		//	authPass = val
		//}
		//if val := params.Get("_auth_crypt"); val != "" {
		//	authCrypt = val
		//}
		//if val := params.Get("_auth_salt"); val != "" {
		//	authSalt = val
		//}

		// COMMENT_FLAG: don't support
		// _loc
		//if val := params.Get("_loc"); val != "" {
		//	switch strings.ToLower(val) {
		//	case "auto":
		//		loc = time.Local
		//	default:
		//		loc, err = time.LoadLocation(val)
		//		if err != nil {
		//			return nil, fmt.Errorf("invalid _loc: %v: %v", val, err)
		//		}
		//	}
		//}

		// COMMENT_FLAG: only sqlite3.SQLITE_OPEN_FULLMUTEX
		// _mutex
		//if val := params.Get("_mutex"); val != "" {
		//	switch strings.ToLower(val) {
		//	case "no":
		//		mutex = C.SQLITE_OPEN_NOMUTEX
		//	case "full":
		//		mutex = C.SQLITE_OPEN_FULLMUTEX
		//	default:
		//		return nil, fmt.Errorf("invalid _mutex: %v", val)
		//	}
		//}

		// COMMENT_FLAG: do in sqlite3.Open()
		// _txlock
		//if val := params.Get("_txlock"); val != "" {
		//	switch strings.ToLower(val) {
		//	case "immediate":
		//		txlock = "BEGIN IMMEDIATE"
		//	case "exclusive":
		//		txlock = "BEGIN EXCLUSIVE"
		//	case "deferred":
		//		txlock = "BEGIN"
		//	default:
		//		return nil, fmt.Errorf("invalid _txlock: %v", val)
		//	}
		//}

		// Auto Vacuum (_vacuum)
		//
		// https://www.sqlite.org/pragma.html#pragma_auto_vacuum
		//
		pkey = "" // Reset pkey
		if _, ok := params["_auto_vacuum"]; ok {
			pkey = "_auto_vacuum"
		}
		if _, ok := params["_vacuum"]; ok {
			pkey = "_vacuum"
		}
		if val := params.Get(pkey); val != "" {
			switch strings.ToLower(val) {
			case "0", "none":
				autoVacuum = 0
			case "1", "full":
				autoVacuum = 1
			case "2", "incremental":
				autoVacuum = 2
			default:
				return nil, fmt.Errorf("invalid _auto_vacuum: %v, expecting value of '0 NONE 1 FULL 2 INCREMENTAL'", val)
			}
		}

		// Busy Timeout (_busy_timeout)
		//
		// https://www.sqlite.org/pragma.html#pragma_busy_timeout
		//
		pkey = "" // Reset pkey
		if _, ok := params["_busy_timeout"]; ok {
			pkey = "_busy_timeout"
		}
		if _, ok := params["_timeout"]; ok {
			pkey = "_timeout"
		}
		if val := params.Get(pkey); val != "" {
			iv, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid _busy_timeout: %v: %v", val, err)
			}
			busyTimeout = int(iv)
		}

		// Case Sensitive Like (_cslike)
		//
		// https://www.sqlite.org/pragma.html#pragma_case_sensitive_like
		//
		pkey = "" // Reset pkey
		if _, ok := params["_case_sensitive_like"]; ok {
			pkey = "_case_sensitive_like"
		}
		if _, ok := params["_cslike"]; ok {
			pkey = "_cslike"
		}
		if val := params.Get(pkey); val != "" {
			switch strings.ToLower(val) {
			case "0", "no", "false", "off":
				caseSensitiveLike = 0
			case "1", "yes", "true", "on":
				caseSensitiveLike = 1
			default:
				return nil, fmt.Errorf("invalid _case_sensitive_like: %v, expecting boolean value of '0 1 false true no yes off on'", val)
			}
		}

		// Defer Foreign Keys (_defer_foreign_keys | _defer_fk)
		//
		// https://www.sqlite.org/pragma.html#pragma_defer_foreign_keys
		//
		pkey = "" // Reset pkey
		if _, ok := params["_defer_foreign_keys"]; ok {
			pkey = "_defer_foreign_keys"
		}
		if _, ok := params["_defer_fk"]; ok {
			pkey = "_defer_fk"
		}
		if val := params.Get(pkey); val != "" {
			switch strings.ToLower(val) {
			case "0", "no", "false", "off":
				deferForeignKeys = 0
			case "1", "yes", "true", "on":
				deferForeignKeys = 1
			default:
				return nil, fmt.Errorf("invalid _defer_foreign_keys: %v, expecting boolean value of '0 1 false true no yes off on'", val)
			}
		}

		// Foreign Keys (_foreign_keys | _fk)
		//
		// https://www.sqlite.org/pragma.html#pragma_foreign_keys
		//
		pkey = "" // Reset pkey
		if _, ok := params["_foreign_keys"]; ok {
			pkey = "_foreign_keys"
		}
		if _, ok := params["_fk"]; ok {
			pkey = "_fk"
		}
		if val := params.Get(pkey); val != "" {
			switch strings.ToLower(val) {
			case "0", "no", "false", "off":
				foreignKeys = 0
			case "1", "yes", "true", "on":
				foreignKeys = 1
			default:
				return nil, fmt.Errorf("invalid _foreign_keys: %v, expecting boolean value of '0 1 false true no yes off on'", val)
			}
		}

		// Ignore CHECK Constrains (_ignore_check_constraints)
		//
		// https://www.sqlite.org/pragma.html#pragma_ignore_check_constraints
		//
		if val := params.Get("_ignore_check_constraints"); val != "" {
			switch strings.ToLower(val) {
			case "0", "no", "false", "off":
				ignoreCheckConstraints = 0
			case "1", "yes", "true", "on":
				ignoreCheckConstraints = 1
			default:
				return nil, fmt.Errorf("invalid _ignore_check_constraints: %v, expecting boolean value of '0 1 false true no yes off on'", val)
			}
		}

		// Journal Mode (_journal_mode | _journal)
		//
		// https://www.sqlite.org/pragma.html#pragma_journal_mode
		//
		pkey = "" // Reset pkey
		if _, ok := params["_journal_mode"]; ok {
			pkey = "_journal_mode"
		}
		if _, ok := params["_journal"]; ok {
			pkey = "_journal"
		}
		if val := params.Get(pkey); val != "" {
			switch strings.ToUpper(val) {
			case "DELETE", "TRUNCATE", "PERSIST", "MEMORY", "OFF":
				journalMode = strings.ToUpper(val)
			case "WAL":
				journalMode = strings.ToUpper(val)

				// For WAL Mode set Synchronous Mode to 'NORMAL'
				// See https://www.sqlite.org/pragma.html#pragma_synchronous
				synchronousMode = "NORMAL"
			default:
				return nil, fmt.Errorf("invalid _journal: %v, expecting value of 'DELETE TRUNCATE PERSIST MEMORY WAL OFF'", val)
			}
		}

		// Locking Mode (_locking)
		//
		// https://www.sqlite.org/pragma.html#pragma_locking_mode
		//
		pkey = "" // Reset pkey
		if _, ok := params["_locking_mode"]; ok {
			pkey = "_locking_mode"
		}
		if _, ok := params["_locking"]; ok {
			pkey = "_locking"
		}
		if val := params.Get(pkey); val != "" {
			switch strings.ToUpper(val) {
			case "NORMAL", "EXCLUSIVE":
				lockingMode = strings.ToUpper(val)
			default:
				return nil, fmt.Errorf("invalid _locking_mode: %v, expecting value of 'NORMAL EXCLUSIVE", val)
			}
		}

		// Query Only (_query_only)
		//
		// https://www.sqlite.org/pragma.html#pragma_query_only
		//
		if val := params.Get("_query_only"); val != "" {
			switch strings.ToLower(val) {
			case "0", "no", "false", "off":
				queryOnly = 0
			case "1", "yes", "true", "on":
				queryOnly = 1
			default:
				return nil, fmt.Errorf("invalid _query_only: %v, expecting boolean value of '0 1 false true no yes off on'", val)
			}
		}

		// Recursive Triggers (_recursive_triggers)
		//
		// https://www.sqlite.org/pragma.html#pragma_recursive_triggers
		//
		pkey = "" // Reset pkey
		if _, ok := params["_recursive_triggers"]; ok {
			pkey = "_recursive_triggers"
		}
		if _, ok := params["_rt"]; ok {
			pkey = "_rt"
		}
		if val := params.Get(pkey); val != "" {
			switch strings.ToLower(val) {
			case "0", "no", "false", "off":
				recursiveTriggers = 0
			case "1", "yes", "true", "on":
				recursiveTriggers = 1
			default:
				return nil, fmt.Errorf("invalid _recursive_triggers: %v, expecting boolean value of '0 1 false true no yes off on'", val)
			}
		}

		// Secure Delete (_secure_delete)
		//
		// https://www.sqlite.org/pragma.html#pragma_secure_delete
		//
		if val := params.Get("_secure_delete"); val != "" {
			switch strings.ToLower(val) {
			case "0", "no", "false", "off":
				secureDelete = "OFF"
			case "1", "yes", "true", "on":
				secureDelete = "ON"
			case "fast":
				secureDelete = "FAST"
			default:
				return nil, fmt.Errorf("invalid _secure_delete: %v, expecting boolean value of '0 1 false true no yes off on fast'", val)
			}
		}

		// Synchronous Mode (_synchronous | _sync)
		//
		// https://www.sqlite.org/pragma.html#pragma_synchronous
		//
		pkey = "" // Reset pkey
		if _, ok := params["_synchronous"]; ok {
			pkey = "_synchronous"
		}
		if _, ok := params["_sync"]; ok {
			pkey = "_sync"
		}
		if val := params.Get(pkey); val != "" {
			switch strings.ToUpper(val) {
			case "0", "OFF", "1", "NORMAL", "2", "FULL", "3", "EXTRA":
				synchronousMode = strings.ToUpper(val)
			default:
				return nil, fmt.Errorf("invalid _synchronous: %v, expecting value of '0 OFF 1 NORMAL 2 FULL 3 EXTRA'", val)
			}
		}

		// Writable Schema (_writeable_schema)
		//
		// https://www.sqlite.org/pragma.html#pragma_writeable_schema
		//
		if val := params.Get("_writable_schema"); val != "" {
			switch strings.ToLower(val) {
			case "0", "no", "false", "off":
				writableSchema = 0
			case "1", "yes", "true", "on":
				writableSchema = 1
			default:
				return nil, fmt.Errorf("invalid _writable_schema: %v, expecting boolean value of '0 1 false true no yes off on'", val)
			}
		}

		// Cache size (_cache_size)
		//
		// https://sqlite.org/pragma.html#pragma_cache_size
		//
		if val := params.Get("_cache_size"); val != "" {
			iv, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid _cache_size: %v: %v", val, err)
			}
			cacheSize = &iv
		}

		//if val := params.Get("vfs"); val != "" {
		//	vfsName = val
		//}
		//
		//if !strings.HasPrefix(dsn, "file:") {
		//	dsn = dsn[:pos]
		//}
	}

	// Open sqlite3 database
	c, err := d.drv.Open(dsn)
	if err != nil {
		return nil, err
	}

	conn, ok := c.(sqliteConn)
	if !ok {
		return c, nil
	}

	exec := func(s string) error {
		_, err := conn.Exec(s, nil)
		return err
	}

	// Busy timeout
	if err := exec(fmt.Sprintf("PRAGMA busy_timeout = %d;", busyTimeout)); err != nil {
		_ = conn.Close()
		return nil, err
	}

	// Auto Vacuum
	// Moved auto_vacuum command, the user preference for auto_vacuum needs to be implemented directly after
	// the authentication and before the sqlite_user table gets created if the user
	// decides to activate User Authentication because
	// auto_vacuum needs to be set before any tables are created
	// and activating user authentication creates the internal table `sqlite_user`.
	if autoVacuum > -1 {
		if err := exec(fmt.Sprintf("PRAGMA auto_vacuum = %d;", autoVacuum)); err != nil {
			_ = conn.Close()
			return nil, err
		}
	}

	// Case Sensitive LIKE
	if caseSensitiveLike > -1 {
		if err := exec(fmt.Sprintf("PRAGMA case_sensitive_like = %d;", caseSensitiveLike)); err != nil {
			_ = conn.Close()
			return nil, err
		}
	}

	// Defer Foreign Keys
	if deferForeignKeys > -1 {
		if err := exec(fmt.Sprintf("PRAGMA defer_foreign_keys = %d;", deferForeignKeys)); err != nil {
			_ = conn.Close()
			return nil, err
		}
	}

	// Foreign Keys
	if foreignKeys > -1 {
		if err := exec(fmt.Sprintf("PRAGMA foreign_keys = %d;", foreignKeys)); err != nil {
			_ = conn.Close()
			return nil, err
		}
	}

	// Ignore CHECK Constraints
	if ignoreCheckConstraints > -1 {
		if err := exec(fmt.Sprintf("PRAGMA ignore_check_constraints = %d;", ignoreCheckConstraints)); err != nil {
			_ = conn.Close()
			return nil, err
		}
	}

	// Journal Mode
	if journalMode != "" {
		if err := exec(fmt.Sprintf("PRAGMA journal_mode = %s;", journalMode)); err != nil {
			_ = conn.Close()
			return nil, err
		}
	}

	// Locking Mode
	// Because the default is NORMAL and this is not changed in this package
	// by using the compile time SQLITE_DEFAULT_LOCKING_MODE this PRAGMA can always be executed
	if err := exec(fmt.Sprintf("PRAGMA locking_mode = %s;", lockingMode)); err != nil {
		_ = conn.Close()
		return nil, err
	}

	// Query Only
	if queryOnly > -1 {
		if err := exec(fmt.Sprintf("PRAGMA query_only = %d;", queryOnly)); err != nil {
			_ = conn.Close()
			return nil, err
		}
	}

	// Recursive Triggers
	if recursiveTriggers > -1 {
		if err := exec(fmt.Sprintf("PRAGMA recursive_triggers = %d;", recursiveTriggers)); err != nil {
			_ = conn.Close()
			return nil, err
		}
	}

	// Secure Delete
	//
	// Because this package can set the compile time flag SQLITE_SECURE_DELETE with a build tag
	// the default value for secureDelete var is 'DEFAULT' this way
	// you can compile with secure_delete 'ON' and disable it for a specific database connection.
	if secureDelete != "DEFAULT" {
		if err := exec(fmt.Sprintf("PRAGMA secure_delete = %s;", secureDelete)); err != nil {
			_ = conn.Close()
			return nil, err
		}
	}

	// Synchronous Mode
	//
	// Because default is NORMAL this statement is always executed
	if err := exec(fmt.Sprintf("PRAGMA synchronous = %s;", synchronousMode)); err != nil {
		_ = conn.Close()
		return nil, err
	}

	// Writable Schema
	if writableSchema > -1 {
		if err := exec(fmt.Sprintf("PRAGMA writable_schema = %d;", writableSchema)); err != nil {
			_ = conn.Close()
			return nil, err
		}
	}

	// Cache Size
	if cacheSize != nil {
		if err := exec(fmt.Sprintf("PRAGMA cache_size = %d;", *cacheSize)); err != nil {
			_ = conn.Close()
			return nil, err
		}
	}

	// COMMENT_FLAG: never use this
	//if len(d.Extensions) > 0 {
	//	if err := conn.loadExtensions(d.Extensions); err != nil {
	//		conn.Close()
	//		return nil, err
	//	}
	//}
	//
	//if d.ConnectHook != nil {
	//	if err := d.ConnectHook(conn); err != nil {
	//		conn.Close()
	//		return nil, err
	//	}
	//}
	runtime.SetFinalizer(conn, sqliteConn.Close)
	return c, nil
}
