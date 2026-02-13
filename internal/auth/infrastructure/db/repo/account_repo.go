package repo

import (
	"context"
	"database/sql"

	"github.com/artcodefun/heat-expansion-server/internal/auth/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/auth/domain"
	"github.com/artcodefun/heat-expansion-server/internal/auth/infrastructure/db/gen"
	"github.com/google/uuid"
)

type AccountRepository struct {
	db *gen.Queries
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{
		db: gen.New(db),
	}
}

func (r *AccountRepository) Tx(tx ports.Transaction) ports.AccountRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &AccountRepository{
			db: r.db.WithTx(sqlTx),
		}
	}
	return r
}

func (r *AccountRepository) Create(ctx context.Context, acc *domain.Account) error {
	row, err := r.db.CreateAccount(ctx, gen.CreateAccountParams{
		ID:           acc.ID,
		Name:         acc.Name,
		Email:        acc.Email,
		PasswordHash: acc.PasswordHash,
	})
	if err != nil {
		return err
	}
	acc.ID = row.ID
	return nil
}

func (r *AccountRepository) FindByEmail(ctx context.Context, email string) (*domain.Account, error) {
	row, err := r.db.GetAccountByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &domain.Account{
		ID:           row.ID,
		Name:         row.Name,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
	}, nil
}

func (r *AccountRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	row, err := r.db.GetAccountByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &domain.Account{
		ID:           row.ID,
		Name:         row.Name,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
	}, nil
}
