package internal

import "time"

const (
	Applied MigrationStatus = iota
	Pending
	Error
)

type Dialect string

type MigrationStatus int

type MigrationDbModel struct {
	Version   int64           `db:"version"`
	Status    MigrationStatus `db:"status"`
	AppliedAt time.Time       `db:"applied_at"`
}

type MigrationListFilter struct {
	Offset int
	Limit  int
}
