package internal

import (
	"context"
)

type ISqlDialect interface {
	CreateVersionTable(ctx context.Context) error
	InsertVersion(ctx context.Context, version int64, status MigrationStatus) error
	RemoveVersion(ctx context.Context, version int64) error
	GetMigrationsList(ctx context.Context, filter *MigrationListFilter) ([]MigrationDbModel, error)
	GetDbVersion(ctx context.Context) (int64, error)
}
