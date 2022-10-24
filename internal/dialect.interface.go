package internal

import (
	"context"
)

type ISqlDialect interface {
	CreateVersionTable(ctx context.Context) error
	InsertVersion(ctx context.Context, version int64) error
	RemoveVersion(ctx context.Context, version int64) error
	GetMigrationsList(ctx context.Context, filter *MigrationListFilter) (MigrationRecords, error)
	GetDbVersion(ctx context.Context) (int64, error)
}
