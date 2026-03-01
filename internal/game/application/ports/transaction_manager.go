package ports

import "context"

// Transaction is an opaque handle representing a database transaction.
// IMPORTANT: Commit/Rollback are intentionally NOT exposed here. The lifecycle
// (begin/commit/rollback) is owned by the TransactionManager.WithTx implementation.
// Repositories should only use this handle to bind themselves to the same DB tx
// (via the Repo.Tx(tx) pattern). Infrastructure will provide the concrete value
// (e.g., *sql.Tx) behind this opaque interface.
type Transaction interface{}

// TransactionManager coordinates execution of a function within a single database transaction.
// Semantics:
// - Starts a DB transaction
// - Invokes fn with the tx handle
// - If fn returns nil: COMMIT
// - If fn returns error or panics: ROLLBACK (re-panics after rollback)
// This keeps commit/rollback out of application code and prevents misuse.
//
// Usage pattern in use cases:
//
//	TxMgr.WithTx(ctx, func(tx Transaction) error {
//	    baseRepo := userBaseRepo.Tx(tx)
//	    if err := baseRepo.Update(base); err != nil { return err }
//	    // enqueue outbox, etc.
//	    return nil
//	})
type TransactionManager interface {
	WithTx(ctx context.Context, fn func(tx Transaction) error) error
}
