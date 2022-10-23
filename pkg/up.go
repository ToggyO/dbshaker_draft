package dbshaker

import (
	"database/sql"
	"github.com/ToggyO/dbshaker/internal"
)

func Up(db *sql.DB, directory string) error {
	if migrator == nil {
		migrator = internal.NewMigrationRunner(db)
	}

}
