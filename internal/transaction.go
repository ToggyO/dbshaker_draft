package internal

import (
	"context"
	"database/sql"
)

const transactionKey = "t_x_transaction"

type TransactionManager struct {
	db *sql.DB
}

func (tm *TransactionManager) Transaction(ctx context.Context, action TransactionAction) error {
	tx, err := tm.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			// TODO: check
			//if err = tx.Rollback(); err != nil {
			//	return err
			//}
		} else if err != nil {
			xerr := tx.Rollback()
			if xerr != nil {
				err = xerr
			}
		} else {
			err = tx.Commit()
		}
	}()

	ctx = context.WithValue(ctx, transactionKey, tx)

	err = action(ctx)
	return err
}

func (tm *TransactionManager) GetQueryRunner(ctx context.Context) IQueryRunner {
	if txRunner, ok := ctx.Value(transactionKey).(*sql.Tx); ok {
		return txRunner
	}
	return tm.db
}
