package internal

import (
	"context"
	"database/sql"
	"fmt"
)

type postgresDialect struct {
	db        *sql.DB
	tableName string
}

func NewPostgresDialect(db *sql.DB, tableName string) ISqlDialect {
	return &postgresDialect{
		db:        db,
		tableName: tableName,
	}
}

// TODO: Подумать, где воткнуть транзакцию
func (p postgresDialect) CreateVersionTable(ctx context.Context) error {
	query := fmt.Sprintf(`CREATE TABLE %s (
			version BIGINT NOT NULL,
			status SMALLINT NOT NULL,
			applied_at DATE DEFAULT NOW()
	);`, p.tableName)

	_, err := p.db.ExecContext(ctx, query)
	return err
}

func (p postgresDialect) InsertVersion(ctx context.Context, version int64, status MigrationStatus) error {
	query := fmt.Sprintf(`INSERT INTO %s (version,status,applied_at) VALUES ($1, $2);`, p.tableName)
	_, err := p.db.ExecContext(ctx, query, version, status)
	return err
}

func (p postgresDialect) RemoveVersion(ctx context.Context, version int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE version = $1;`, p.tableName)
	_, err := p.db.ExecContext(ctx, query, version)
	return err
}

func (p postgresDialect) GetMigrationsList(ctx context.Context, filter *MigrationListFilter) ([]MigrationDbModel, error) {
	query := fmt.Sprintf(`SELECT version, status, applied_at FROM %s OFFSET $1 LIMIT $2`, p.tableName)

	var offset int
	var limit int
	if filter != nil {
		offset = filter.Offset
		limit = filter.Limit
	}

	rows, err := p.db.QueryContext(ctx, query, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query migrations: %w", err)
	}
	var migrations []MigrationDbModel

	for rows.Next() {
		var model MigrationDbModel

		if err := rows.Scan(&model.Version, &model.Status, &model.AppliedAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		migrations = append(migrations, model)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get next row: %w", err)
	}

	return migrations, nil
}

func (p postgresDialect) GetDbVersion(ctx context.Context) (int64, error) {
	query := fmt.Sprintf(`SELECT version FROM %s ORDER BY version DESC;`, p.tableName)

	var version int64
	if err := p.db.QueryRowContext(ctx, query).Scan(&version); err != nil {
		return 0, err
	}

	return version, nil
}
