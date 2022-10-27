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
	migrator.setDb(db)
	dialect := migrator.getDialect()

	migrations, err := lookupMigrations(directory)
	if err != nil {
		return err
	}

	migrationsMap := make(map[int64]*internal.Migration)
	for _, m := range migrations {
		migrationsMap[m.Version] = m
	}

	for {
		currentDbVersion, err := EnsureDbVersionContext(ctx, db)
		if err != nil {
			return err
		}

		if currentDbVersion == 0 {
			internal.LogWithPrefix(fmt.Sprintf("no migrations to run. current version: %d\n", currentDbVersion))
			return nil
		}

		currentMigration, ok := migrationsMap[currentDbVersion]
		if !ok {
			internal.LogWithPrefix(fmt.Sprintf("goose: no migrations to run. current version: %d\n", currentDbVersion))
			return nil
		}

		if currentMigration.Version < targetVersion {
			internal.LogWithPrefix(fmt.Sprintf("goose: no migrations to run. current version: %d\n", currentDbVersion))
			return nil
		}

		if err = currentMigration.DownContext(ctx, db, dialect); err != nil {
			return err
		}
	}
}
