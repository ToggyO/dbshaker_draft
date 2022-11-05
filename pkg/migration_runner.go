package dbshaker

import (
	"database/sql"

	"github.com/ToggyO/dbshaker/internal"
)

var migrator = &migrationRunner{
	registeredGoMigrations: make(map[int64]*internal.Migration),
}

type migrationRunner struct {
	db                     *sql.DB
	dialect                internal.ISqlDialect
	registeredGoMigrations map[int64]*internal.Migration
}

func newMigrationRunner(db *DB) *migrationRunner {
	return &migrationRunner{
		db:                     db.db,
		dialect:                db.dialect,
		registeredGoMigrations: make(map[int64]*internal.Migration),
	}
}

func (mr *migrationRunner) setDb(db *sql.DB) {
	mr.db = db
}

func (mr *migrationRunner) getDb() *sql.DB {
	return mr.db
}

func (mr *migrationRunner) setDialect(dialect internal.ISqlDialect) {
	mr.dialect = dialect
}

func (mr *migrationRunner) getDialect() internal.ISqlDialect {
	return mr.dialect
}
