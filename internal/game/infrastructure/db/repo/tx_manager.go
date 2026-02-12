package repo

import (
	"context"
	"database/sql"

	"github.com/artcodefun/heat-expansion-api/internal/game/core/ports"
)

// DBTxManager is a minimal TransactionManager over *sql.DB.
// It includes a CommitSignal channel that background workers can use
// to be notified immediately after a successful transaction commit.
type DBTxManager struct {
	db           *sql.DB
	CommitSignal chan struct{}
}

func NewDBTxManager(db *sql.DB) *DBTxManager {
	return &DBTxManager{
		db:           db,
		CommitSignal: make(chan struct{}, 1),
	}
}

func (m *DBTxManager) WithTx(fn func(tx ports.Transaction) error) error {
	tx, err := m.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	// Ensure rollback on panic or error
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	// Non-blocking signal
	select {
	case m.CommitSignal <- struct{}{}:
	default:
	}

	return nil
}
