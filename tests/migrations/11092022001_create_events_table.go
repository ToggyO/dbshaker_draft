package migrations

import (
	"database/sql"

	"github.com/ToggyO/dbshaker/pkg"
)

func init() {
	dbshaker.AddMigration(Up11092022001, Down11092022001)
}

func Up11092022001(tx *sql.Tx) error {
	_, err := tx.Exec(
		`CREATE TABLE events(
		id SERIAL PRIMARY KEY,
		title VARCHAR NOT NULL,
		start_date TIMESTAMP NOT NULL,
		end_date TIMESTAMP,
		description TEXT,
		owner_id INTEGER NOT NULL,
		notification_date TIMESTAMP NOT NULL,
		created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP
   	);`)
	if err != nil {
		return err
	}
	return nil
}

func Down11092022001(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE events;")
	if err != nil {
		return err
	}
	return nil
}
