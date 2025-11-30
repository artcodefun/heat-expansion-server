package repo

import (
	"context"
	"database/sql"

	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
)

// dbTxManager is a minimal TransactionManager over *sql.DB.
type dbTxManager struct {
	db *sql.DB
}

func NewDBTxManager(db *sql.DB) ports.TransactionManager {
	return &dbTxManager{db: db}
}

func (m *dbTxManager) WithTx(fn func(tx ports.Transaction) error) error {
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
	return tx.Commit()
}
