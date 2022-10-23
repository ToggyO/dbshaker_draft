package internal

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
)

type MigrationFunc func(tx *sql.Tx) error

type Migration struct {
	Name    string
	Version int64

	UpFn   MigrationFunc
	DownFn MigrationFunc

	IsApplied bool
}

func (m *Migration) Up(db *sql.DB) error {
	return m.run(db, true)
}

func (m *Migration) Down(db *sql.DB) error {
	return m.run(db, false)
}

func (m *Migration) run(db *sql.DB, direction bool) error {
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
				tx.Rollback()
				return fmt.Errorf("ERROR %v: failed to run Go migration function %T: %w", filepath.Base(m.Name), fn, err)
			}

		}

		// TODO: insert version

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
