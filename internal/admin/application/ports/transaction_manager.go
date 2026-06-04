package ports

import "context"

// Transaction is an opaque handle for a database transaction.
// Repositories use it to bind themselves to the active transaction via Repo.Tx(tx).
type Transaction interface{}

// TransactionManager runs a function within a single database transaction.
// The transaction is committed when fn returns nil and rolled back on error.
type TransactionManager interface {
	WithTx(ctx context.Context, fn func(tx Transaction) error) error
}
