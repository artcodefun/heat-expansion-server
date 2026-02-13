package ports

// Transaction is an opaque handle representing a database transaction.
type Transaction interface{}

// TransactionManager coordinates execution of a function within a single database transaction.
type TransactionManager interface {
	WithTx(fn func(tx Transaction) error) error
}
