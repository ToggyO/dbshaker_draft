package dbshaker

import (
	"fmt"
	"runtime"

	"github.com/ToggyO/dbshaker/internal"
)

var (
	registeredGoMigrations = make(map[int64]*internal.Migration)
)

type Migrations []*internal.Migration

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

	if exists, ok := migrator.TryGetMigration(version); ok {
		panic(fmt.Sprintf("failed to add migration %q: conflicts with exitsting %q", filename, exists.Name))
	}

	registeredGoMigrations[version] = migration
}

func LookupMigrations(directory string) ([]internal.Migration, error) {
	var migrations Migrations

	// Migrations in .go files, registered via AddMigration
	for k, v := registeredGoMigrations {

	}
}
