package dbshaker

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ToggyO/dbshaker/internal"
)

// Up - migrates up to a max version.
func Up(db *sql.DB, directory string) error {
	return UpTo(db, directory, maxVersion)
}

// UpContext migrates up to a max version with context.
func UpContext(ctx context.Context, db *sql.DB, directory string) error {
	return UpToContext(ctx, db, directory, maxVersion)
}

// UpTo migrates up to a specific version.
func UpTo(db *sql.DB, directory string, targetVersion int64) error {
	return UpToContext(context.Background(), db, directory, targetVersion)
}

// UpToContext migrates up to a specific version with context.
func UpToContext(ctx context.Context, db *sql.DB, directory string, targetVersion int64) error {
	migrator.setDb(db)
	dialect := migrator.getDialect()

	// TODO: check
	//return dialect.Transaction(ctx, func(ctx context.Context) error {
	//	foundMigrations, err := lookupMigrations(directory, targetVersion)
	//	if err != nil {
	//		return err
	//	}
	//
	//	currentDbVersion, err := EnsureDbVersionContext(ctx, db)
	//	if err != nil {
	//		return err
	//	}
	//
	//	dbMigrations, err := dialect.GetMigrationsList(ctx, nil)
	//	if err != nil {
	//		return err
	//	}
	//
	//	notAppliedMigrations := lookupNotAppliedMigrations(dbMigrations.ToMigrationsList(), foundMigrations)
	//
	//	for _, migration := range notAppliedMigrations {
	//		if err = migration.UpContext(ctx, db, dialect); err != nil {
	//			return err
	//		}
	//	}
	//
	//	notAppliedMigrationsLen := len(notAppliedMigrations)
	//	if notAppliedMigrationsLen > 0 {
	//		if notAppliedMigrations[notAppliedMigrationsLen-1].Version < currentDbVersion {
	//			err = dialect.IncrementVersionPatch(ctx, currentDbVersion)
	//		}
	//	}
	//
	//	internal.LogWithPrefix(fmt.Sprintf("no migrations to run. current version: %d\n", currentDbVersion))
	//	return nil
	//})

	foundMigrations, err := lookupMigrations(directory, targetVersion)
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

	notAppliedMigrationsLen := len(notAppliedMigrations)
	if notAppliedMigrationsLen > 0 {
		if notAppliedMigrations[notAppliedMigrationsLen-1].Version < currentDbVersion {
			err = dialect.IncrementVersionPatch(ctx, currentDbVersion)
		}
	}

	internal.LogWithPrefix(fmt.Sprintf("no migrations to run. current version: %d\n", currentDbVersion))
	return nil
}
