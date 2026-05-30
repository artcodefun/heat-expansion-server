package ports

import "context"

type Transaction interface{}

type TransactionManager interface {
	WithTx(ctx context.Context, fn func(tx Transaction) error) error
}
