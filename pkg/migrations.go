package dbshaker

import (
	"github.com/ToggyO/dbshaker/internal"
	"path/filepath"
	"sort"
)

const MaxUint = ^uint64(0)
const maxVersion = int64(MaxUint >> 1) // max(int64)

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
		panic(internal.ErrDuplicateVersion(ms[i].Version, ms[i].Source, ms[j].Source))
	}
	return ms[i].Version < ms[j].Version
}

// LookupMigrations returns a slice of valid migrations in the migrations folder and migration registry,
// sorted by version in ascending direction.
// TODO: `embed` support in future by embed.FS
func lookupMigrations(directory string, targetVersion int64) (Migrations, error) {
	var migrations Migrations

	// SQL migrations
	//sqlMigrationFiles, err := fs.Glob() for embedding `.sql` migrations
	sqlMigrationFiles, err := filepath.Glob(directory + internal.SqlFilesPattern)
	if err != nil {
		return nil, err
	}

	// micro optimization of migrations slice allocation size
	if len(sqlMigrationFiles) > 0 {
		migrations = make(Migrations, 0, len(sqlMigrationFiles)+len(migrator.registeredGoMigrations))
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
		migrations = make(Migrations, 0, len(migrator.registeredGoMigrations))
	}

	// Migrations in `.go` files, registered via AddMigration
	for _, migration := range migrator.registeredGoMigrations {
		if migration.Version > targetVersion {
			continue
		}
		migrations = append(migrations, migration)
	}

	// Unregistered `.go` migrations
	unregisteredGoMigrations, err := filepath.Glob(directory + internal.GoFilesPattern)
	if err != nil {
		return nil, err
	}

	if len(unregisteredGoMigrations) > 0 {
		return nil, internal.ErrUnregisteredGoMigration
	}

	// TODO: check sort stability
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

	// TODO: check sort stability
	sort.Sort(migrations)
	return migrations
}
