package dbshaker

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ToggyO/dbshaker/internal"
)

type DB struct {
	db      *sql.DB
	dialect internal.ISqlDialect
}

func OpenDbWithDriver(dialect, connectionString string) (*sql.DB, error) {
	fmt.Printf("Connecting to `%s` database...", dialect)

	var db *sql.DB
	var err error

	switch dialect {
	// tODO: check
	//case "postgres", "pgx", "sqlite3", "sqlite", "mysql", "sqlserver":
	case internal.PostgresDialect, internal.PgxDialect:
		db, err = sql.Open(dialect, connectionString)
	default:
		return nil, fmt.Errorf("unsupported driver '%s'", dialect)
	}

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ERROR: failed connect to database: %v", err)
	}

	if err = SetDialect(db, dialect); err != nil {
		return nil, err
	}

	fmt.Println("Connected to database!")

	return db, nil
}

// EnsureDbVersion retrieves the current version for this DB (major version, patch).
// Create and initialize the DB version table if it doesn't exist.
func EnsureDbVersion(db *sql.DB) (int64, byte, error) {
	return EnsureDbVersionContext(context.Background(), db)
}

// EnsureDbVersionContext retrieves the current version for this DB (major version, patch) with context.
// Create and initialize the DB version table if it doesn't exist.
func EnsureDbVersionContext(ctx context.Context, db *sql.DB) (int64, byte, error) {
	sqlDialect := migrator.getDialect()

	version, err := sqlDialect.GetDbVersion(ctx)
	if err != nil {
		return version.Version, version.Patch, sqlDialect.CreateVersionTable(ctx)
	}

	return version.Version, version.Patch, nil
}
