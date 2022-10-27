package internal

import "context"

// TransactionAction function that will be executed while the transaction is running
type TransactionAction = func(ctx context.Context) error

// ITransactionBuilder represent an SQL transaction process runner
type ITransactionBuilder interface {
	Transaction(ctx context.Context, action TransactionAction) error
	GetQueryRunner(ctx context.Context) IQueryRunner
}
