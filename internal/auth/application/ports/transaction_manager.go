package ports

import "context"

// Transaction is an opaque handle representing a database transaction.
type Transaction interface{}

// TransactionManager coordinates execution of a function within a single database transaction.
type TransactionManager interface {
	WithTx(ctx context.Context, fn func(tx Transaction) error) error
}
