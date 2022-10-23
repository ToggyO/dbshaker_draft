package internal

import (
	"database/sql"
)

// MigrationRunner - deprecated
type MigrationRunner struct {
	db                   *sql.DB
	registeredMigrations map[int64]*Migration
}

func NewMigrationRunner(db *sql.DB) *MigrationRunner {
	return &MigrationRunner{
		db:                   db,
		registeredMigrations: make(map[int64]*Migration),
	}
}

func (mr *MigrationRunner) RegisterMigration(version int64, migration *Migration) {
	mr.registeredMigrations[version] = migration
}

func (mr *MigrationRunner) TryGetMigration(version int64) (*Migration, bool) {
	exists, ok := mr.registeredMigrations[version]
	return exists, ok
}

func (mr *MigrationRunner) Iterator() {

}
