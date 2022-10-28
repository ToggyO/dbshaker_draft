package dbshaker

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ToggyO/dbshaker/internal"
)

// Down rolls back all existing migrations.
func Down(db *sql.DB, directory string) error {
	return DownContext(context.Background(), db, directory)
}

// DownContext rolls back all existing migrations with context.
func DownContext(ctx context.Context, db *sql.DB, directory string) error {
	return DownToContext(ctx, db, directory, 0)
}

// DownTo rolls back migrations to a specific version.
func DownTo(db *sql.DB, directory string, targetVersion int64) error {
	return DownToContext(context.Background(), db, directory, targetVersion)
}

// DownToContext rolls back migrations to a specific version with context.
func DownToContext(ctx context.Context, db *sql.DB, directory string, targetVersion int64) error {
	dialect := migrator.getDialect()
	currentDbVersion, err := EnsureDbVersionContext(ctx, db)
	if err != nil {
		return err
	}

	if currentDbVersion < targetVersion {
		return internal.ErrDbAlreadyIsUpToDate(currentDbVersion)
	}

	return dialect.Transaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		migrations, err := lookupMigrations(directory, maxVersion)
		if err != nil {
			return err
		}

		migrationsMap := make(map[int64]*internal.Migration)
		for _, m := range migrations {
			migrationsMap[m.Version] = m
		}

		for {
			currentDbVersion, err = EnsureDbVersionContext(ctx, db)
			if err != nil {
				return err
			}

			if currentDbVersion == 0 {
				// TODO: duplicate
				internal.LogWithPrefix(fmt.Sprintf("no migrations to run. current version: %d\n", currentDbVersion))
				return nil
			}

			currentMigration, ok := migrationsMap[currentDbVersion]
			if !ok {
				// TODO: duplicate
				internal.LogWithPrefix(fmt.Sprintf("no migrations to run. current version: %d\n", currentDbVersion))
				return nil
			}

			if currentMigration.Version < targetVersion {
				// TODO: duplicate
				internal.LogWithPrefix(fmt.Sprintf("no migrations to run. current version: %d\n", currentDbVersion))
				return nil
			}

			if err = currentMigration.DownContext(ctx, tx, dialect); err != nil {
				return err
			}
		}
	})
}
