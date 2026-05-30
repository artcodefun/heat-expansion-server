package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/auth/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/auth/domain"
	"github.com/artcodefun/heat-expansion-server/internal/auth/infrastructure/db/gen"
	"github.com/google/uuid"
)

type AccountRepository struct {
	q *gen.Queries
}

func NewAccountRepository(q *gen.Queries) *AccountRepository {
	return &AccountRepository{q: q}
}

func (r *AccountRepository) Tx(tx ports.Transaction) ports.AccountRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &AccountRepository{q: r.q.WithTx(sqlTx)}
	}
	return r
}

func (r *AccountRepository) Create(ctx context.Context, acc *domain.Account) error {
	row, err := r.q.CreateAccount(ctx, gen.CreateAccountParams{
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
	row, err := r.q.GetAccountByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
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
	row, err := r.q.GetAccountByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
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

func (r *AccountRepository) UpdatePassword(ctx context.Context, id uuid.UUID, newHash string) error {
	return r.q.UpdateAccountPasswordHash(ctx, gen.UpdateAccountPasswordHashParams{
		ID:           id,
		PasswordHash: newHash,
	})
}
