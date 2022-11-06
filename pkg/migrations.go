package dbshaker

import (
	"context"
	"github.com/ToggyO/dbshaker/internal"
	"path/filepath"
	"sort"
)

const (
	maxUint    = ^uint64(0)
	maxVersion = int64(maxUint >> 1) // max(int64)

)

// Migrations runtime slice of Migration struct pointers.
type Migrations []*internal.Migration

func (ms Migrations) Len() int {
	return len(ms)
}

func (ms Migrations) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}

func (ms Migrations) Less(i, j int) bool {
	if ms[i].Version == ms[j].Version {
		logger.Fatal(internal.ErrDuplicateVersion(ms[i].Version, ms[i].Source, ms[j].Source))
	}
	return ms[i].Version < ms[j].Version
}

// ListMigrations lists all applied migrations in database.
func ListMigrations(db *DB) (Migrations, error) {
	return ListMigrationsContext(context.Background(), db)
}

// ListMigrationsContext lists all applied migrations in database with context.
func ListMigrationsContext(ctx context.Context, db *DB) (Migrations, error) {
	records, err := db.dialect.GetMigrationsList(ctx, nil)
	if err != nil {
		return nil, err
	}
	return records.ToMigrationsList(), nil
}

// LookupMigrations returns a slice of valid migrations in the migrations folder and migration registry,
// sorted by version in ascending direction.
// TODO: `embed` support in future by embed.FS
func lookupMigrations(directory string, targetVersion int64) (Migrations, error) {
	key, err := filepath.Abs(directory)
	if err != nil {
		return nil, err
	}

	folderRegistry, ok := registry[key]
	if !ok {
		folderRegistry = make(folderGoMigrationRegistry)
	}

	var migrations Migrations

	// SQL migrations
	//sqlMigrationFiles, err := fs.Glob() for embedding `.sql` migrations
	sqlMigrationFiles, err := filepath.Glob(filepath.Join(directory, internal.SqlFilesPattern))
	if err != nil {
		return nil, err
	}

	// micro optimization of migrations slice allocation size
	if len(sqlMigrationFiles) > 0 {
		migrations = make(Migrations, 0, len(sqlMigrationFiles)+len(folderRegistry))
	}

	for _, file := range sqlMigrationFiles {
		v, err := internal.IsValidFileName(file)
		if err != nil {
			return nil, internal.ErrCouldNotParseMigration(file, err)
		}

		if v > targetVersion {
			continue
		}

		migrations = append(migrations, &internal.Migration{
			Name:    filepath.Base(file),
			Version: v,
			Source:  file,
		})
	}

	// micro optimization of migrations slice allocation size
	if cap(migrations) <= 0 {
		migrations = make(Migrations, 0, len(folderRegistry))
	}

	// Migrations in `.go` files, registered via AddMigration
	for _, migration := range folderRegistry {
		if migration.Version > targetVersion {
			continue
		}
		migrations = append(migrations, migration)
	}

	// Unregistered `.go` migrations
	gGoMigrationsFiles, err := filepath.Glob(filepath.Join(directory, internal.GoFilesPattern))
	if err != nil {
		return nil, err
	}

	for _, file := range gGoMigrationsFiles {
		v, err := internal.IsValidFileName(file)
		if err != nil {
			continue // Пропускаем файлы, которые не имею версионного префикса
		}

		if _, ok := folderRegistry[v]; !ok {
			return nil, internal.ErrUnregisteredGoMigration
		}
	}

	sort.Sort(migrations)

	return migrations, nil
}

func lookupNotAppliedMigrations(known, found Migrations) Migrations {
	existing := make(map[int64]bool)
	for _, k := range known {
		existing[k.Version] = true
	}

	var migrations Migrations
	for _, f := range found {
		if _, ok := existing[f.Version]; !ok {
			migrations = append(migrations, f)
		}
	}

	sort.Sort(migrations)
	return migrations
}
