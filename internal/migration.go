package internal

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
)

// MigrationFunc migration action in database.
type MigrationFunc func(tx *sql.Tx) error

// Migration represents a database migration, manages by go runtime.
type Migration struct {
	Name    string // migration file name.
	Version int64  // version of migration.

	UpFn   MigrationFunc // Up migrations function.
	DownFn MigrationFunc // Down migrations function.

	Source    string // path to migration file.
	IsApplied bool   // indicates, whether migration is applied to database schema.
}

// Up executes an up migration.
func (m *Migration) Up(db *sql.DB, dialect ISqlDialect) error {
	return m.UpContext(context.Background(), db, dialect)
}

// UpContext executes an up migration with context.
func (m *Migration) UpContext(ctx context.Context, db *sql.DB, dialect ISqlDialect) error {
	return m.run(ctx, db, dialect, true)
}

// Down executes an up migration.
func (m *Migration) Down(db *sql.DB, dialect ISqlDialect) error {
	return m.DownContext(context.Background(), db, dialect)
}

// DownContext executes an up migration with context.
func (m *Migration) DownContext(ctx context.Context, db *sql.DB, dialect ISqlDialect) error {
	return m.run(ctx, db, dialect, true)
}

func (m *Migration) run(ctx context.Context, db *sql.DB, dialect ISqlDialect, direction bool) error {
	ext := filepath.Ext(m.Name)
	switch ext {
	case SqlExt:
	case GoExt:
		// TODO: check
		if m.IsApplied {
			return nil
		}

		tx, err := db.Begin()
		if err != nil {
			return err
		}

		fn := m.UpFn
		if !direction {
			fn = m.DownFn
		}

		if fn != nil {
			if err = fn(tx); err != nil {
				_ = tx.Rollback()
				// TODO: duplicate
				return fmt.Errorf("ERROR %v: failed to run Go migration function %T: %w", filepath.Base(m.Name), fn, err)
			}

		}

		if err = dialect.InsertVersion(ctx, m.Version); err != nil {
			_ = tx.Rollback()
			// TODO: duplicate
			return fmt.Errorf("ERROR %v: failed to run Go migration function %T: %w", filepath.Base(m.Name), fn, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("ERROR failed to commit transaction: %w", err)
		}

		if fn != nil {
			log.Println("OK   ", filepath.Base(m.Name))
		} else {
			log.Println("EMPTY", filepath.Base(m.Name))
		}

		return nil
	}

	return nil
}
