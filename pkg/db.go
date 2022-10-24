package dbshaker

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ToggyO/dbshaker/internal"
)

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
		return nil, fmt.Errorf("unsupported core.go '%s'", dialect)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("unsupported core.go '%s'", dialect)
	}

	if err = SetDialect(db, dialect); err != nil {
		return nil, err
	}

	fmt.Println("Connected to database!")

	return db, nil
}

func EnsureDbVersion(db *sql.DB) (int64, error) {
	return EnsureDbVersionContext(context.Background(), db)
}

func EnsureDbVersionContext(ctx context.Context, db *sql.DB) (int64, error) {
	sqlDialect := migrator.getDialect()

	version, err := sqlDialect.GetDbVersion(ctx)
	if err != nil {
		return 0, sqlDialect.CreateVersionTable(ctx)
	}

	return version, nil
}
