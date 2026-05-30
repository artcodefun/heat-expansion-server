package repo

import (
	"context"
	"database/sql"

	"github.com/artcodefun/heat-expansion-server/internal/billing/application/ports"
)

// DBTxManager is a minimal TransactionManager over *sql.DB.
type DBTxManager struct {
	db *sql.DB
}

func NewDBTxManager(db *sql.DB) *DBTxManager {
	return &DBTxManager{
		db: db,
	}
}

func (m *DBTxManager) WithTx(ctx context.Context, fn func(tx ports.Transaction) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
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

	return nil
}
