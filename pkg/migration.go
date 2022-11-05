package dbshaker

import (
	"database/sql"
	"fmt"
	"github.com/ToggyO/dbshaker/internal"
	"runtime"
)

type IMigration interface {
	Up(tx *sql.Tx) error
	Down(tx *sql.Tx) error
}

type GoMigrationSource struct {
	registeredGoMigrations map[int64]*internal.Migration
}

func (s *GoMigrationSource) AddMigration()

func AddMigration(up internal.MigrationFunc, down internal.MigrationFunc) {
	_, filename, _, _ := runtime.Caller(1)

	version, err := internal.IsValidFileName(filename)
	if err != nil {
		panic(err)
	}

	migration := &internal.Migration{
		Name:    filename,
		Version: version,
		UpFn:    up,
		DownFn:  down,
	}

	if exists, ok := migrator.registeredGoMigrations[version]; ok {
		logger.Fatal(fmt.Sprintf("[dbshaker]: failed to add migration %q: conflicts with exitsting %q", filename, exists.Name))
	}

	migrator.registeredGoMigrations[version] = migration
}
