package dbshaker

import (
	"fmt"
	"github.com/ToggyO/dbshaker/internal"
	"runtime"
)

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
		panic(fmt.Sprintf("failed to add migration %q: conflicts with exitsting %q", filename, exists.Name))
	}

	migrator.registeredGoMigrations[version] = migration
}
