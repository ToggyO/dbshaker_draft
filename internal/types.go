package internal

import "time"

type Dialect string

type MigrationListFilter struct {
	Offset int
	Limit  int
}

type MigrationRecord struct {
	Version   int64     `db:"version"`
	Patch     int32     `db:"patch"`
	AppliedAt time.Time `db:"applied_at"`
}

type MigrationRecords []MigrationRecord

func (mr MigrationRecords) ToMigrationsList() []*Migration {
	migrations := make([]*Migration, 0, len(mr))

	for _, migrationRecord := range mr {
		migrations = append(migrations, &Migration{
			Version: migrationRecord.Version,
		})
	}

	return migrations
}
