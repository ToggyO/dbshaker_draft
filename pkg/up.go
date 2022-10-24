package dbshaker

import (
	"context"
	"database/sql"
	"log"
)

func Up(db *sql.DB, directory string) error {
	return UpContext(context.Background(), db, directory)
}

func UpContext(ctx context.Context, db *sql.DB, directory string) error {
	migrator.setDb(db)
	dialect := migrator.getDialect()

	foundMigrations, err := lookupMigrations(directory)
	if err != nil {
		return err
	}

	currentDbVersion, err := EnsureDbVersionContext(ctx, db)
	if err != nil {
		return err
	}

	dbMigrations, err := dialect.GetMigrationsList(ctx, nil)
	if err != nil {
		return err
	}

	notAppliedMigrations := lookupNotAppliedMigrations(dbMigrations.ToMigrationsList(), foundMigrations)
	for _, migration := range notAppliedMigrations {
		if err = migration.UpContext(ctx, db, dialect); err != nil {
			return err
		}
	}

	log.Printf("dbshaker: no migrations to run. current version: %d\n", currentDbVersion)
	return nil
}
