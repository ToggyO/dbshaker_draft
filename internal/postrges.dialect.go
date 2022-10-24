package internal

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
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
// TODO: Продумать патч версии
func (p *postgresDialect) CreateVersionTable(ctx context.Context) error {
	query := fmt.Sprintf(`CREATE TABLE %s (
			version BIGINT NOT NULL,
			applied_at DATE DEFAULT NOW()
	);`, p.tableName)

	_, err := p.db.ExecContext(ctx, query)
	return err
}

func (p *postgresDialect) InsertVersion(ctx context.Context, version int64) error {
	query := fmt.Sprintf(`INSERT INTO %s (version,patch,applied_at) VALUES ($1,);`, p.tableName)
	_, err := p.db.ExecContext(ctx, query, version)
	return err
}

func (p *postgresDialect) RemoveVersion(ctx context.Context, version int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE version = $1;`, p.tableName)
	_, err := p.db.ExecContext(ctx, query, version)
	return err
}

func (p *postgresDialect) GetMigrationsList(ctx context.Context, filter *MigrationListFilter) (MigrationRecords, error) {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf(`SELECT version, applied_at FROM %s OFFSET $1`, p.tableName))

	var offset int
	var limit int
	if filter != nil {
		offset = filter.Offset
		limit = filter.Limit
	}

	if limit == 0 {
		sb.WriteString("LIMIT ALL")
	} else {
		sb.WriteString("LIMIT $2")
	}

	sb.WriteString(";")

	rows, err := p.db.QueryContext(ctx, sb.String(), offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query migrations: %w", err)
	}
	var migrations MigrationRecords

	for rows.Next() {
		var model MigrationRecord

		if err := rows.Scan(&model.Version, &model.AppliedAt); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		migrations = append(migrations, model)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get next row: %w", err)
	}

	return migrations, nil
}

func (p *postgresDialect) GetDbVersion(ctx context.Context) (int64, error) {
	query := fmt.Sprintf(`SELECT version FROM %s ORDER BY version DESC;`, p.tableName)

	var version int64
	if err := p.db.QueryRowContext(ctx, query).Scan(&version); err != nil {
		return 0, err
	}

	return version, nil
}

// TODO: парметр функцию ввынести в типы
func (p *postgresDialect) Transaction(ctx context.Context, action func(ctx context.Context) error) error {

}
