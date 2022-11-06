package dbshaker

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/ToggyO/dbshaker/internal"
)

func AddMigration(up internal.MigrationFunc, down internal.MigrationFunc) {
	_, filename, _, _ := runtime.Caller(1)

	version, err := internal.IsValidFileName(filename)
	if err != nil {
		panic(err)
	}

	key := filepath.Dir(filename)
	if err != nil {
		panic(err)
	}

	folderRegistry, ok := registry[key]
	if !ok {
		folderRegistry = make(folderGoMigrationRegistry)
	}

	migration := &internal.Migration{
		Name:    filename,
		Version: version,
		UpFn:    up,
		DownFn:  down,
	}

	if exists, ok := folderRegistry[version]; ok {
		logger.Fatal(fmt.Sprintf("failed to add migration %q: conflicts with exitsting %q", filename, exists.Name))
	}

	folderRegistry[version] = migration
	registry[key] = folderRegistry
}

//// IMigrationSource describes developer's go migration source.
//type IMigrationSource interface {
//	Up(tx *sql.Tx) error
//	Down(tx *sql.Tx) error
//}
//
//// GoMigrationSourceRegistry is a registry for developer's go migration sources.
//type GoMigrationSourceRegistry struct {
//	registeredGoMigrations map[int64]*internal.Migration
//}
//
//// NewGoMigrationSourceRegistry creates new instance of GoMigrationSourceRegistry.
//func NewGoMigrationSourceRegistry() *GoMigrationSourceRegistry {
//	return &GoMigrationSourceRegistry{
//		registeredGoMigrations: make(map[int64]*internal.Migration),
//	}
//}
//
//// AddMigrations adds a set of migrations sources to the registry.
//func (s *GoMigrationSourceRegistry) AddMigrations(migrationSources ...IMigrationSource) {
//	for _, ms := range migrationSources {
//		s.AddMigration(ms)
//	}
//}
//
//// AddMigration adds new migrations source to the registry.
//func (s *GoMigrationSourceRegistry) AddMigration(migrationSource IMigrationSource) {
//	migrationName := reflect.TypeOf(migrationSource).Name()
//
//	version, err := internal.IsValidMigrationName(migrationName)
//	if err != nil {
//		panic(err)
//	}
//
//	migration := &internal.Migration{
//		Name:    migrationName,
//		Version: version,
//		UpFn:    migrationSource.Up,
//		DownFn:  migrationSource.Down,
//	}
//
//	if exists, ok := s.registeredGoMigrations[version]; ok {
//		logger.Fatal(fmt.Sprintf("[dbshaker]: failed to add migration %q: conflicts with exitsting %q", migrationName, exists.Name))
//	}
//
//	s.registeredGoMigrations[version] = migration
//}

// TODO: remove
//func AddMigration(up internal.MigrationFunc, down internal.MigrationFunc) {
//	_, filename, _, _ := runtime.Caller(1)
//
//	version, err := internal.IsValidFileName(filename)
//	if err != nil {
//		panic(err)
//	}
//
//	migration := &internal.Migration{
//		Name:    filename,
//		Version: version,
//		UpFn:    up,
//		DownFn:  down,
//	}
//
//	if exists, ok := migrator.registeredGoMigrations[version]; ok {
//		logger.Fatal(fmt.Sprintf("[dbshaker]: failed to add migration %q: conflicts with exitsting %q", filename, exists.Name))
//	}
//
//	migrator.registeredGoMigrations[version] = migration
//}
