package dbshaker

import (
	"database/sql"
	"fmt"
	"github.com/ToggyO/dbshaker/internal"
)

var dialect internal.ISqlDialect

func SetDialect(db *sql.DB, d string) error {
	// TODO: добавить поддержку диалектов других СУБД
	switch d {
	case internal.PostgresDialect, internal.PgxDialect:
		dialect = internal.NewPostgresDialect(db, internal.ServiceTableName)
	default:
		return fmt.Errorf("%q: unknown dialect", d)
	}

	return nil
}
