package internal

import (
	"context"
)

type ISqlDialect interface {
	ITransactionBuilder

	CreateVersionTable(ctx context.Context) error
	InsertVersion(ctx context.Context, version int64) error
	IncrementVersionPatch(ctx context.Context, version int64) error
	RemoveVersion(ctx context.Context, version int64) error
	GetMigrationsList(ctx context.Context, filter *MigrationListFilter) (MigrationRecords, error)
	GetDbVersion(ctx context.Context) (DbVersion, error)
}
