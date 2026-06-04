package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/admin/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/admin/domain"
	"github.com/artcodefun/heat-expansion-server/internal/admin/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-server/internal/admin/infrastructure/db/mappers"
	"github.com/google/uuid"
)

type AdminRepository struct {
	q *gen.Queries
}

func NewAdminRepository(q *gen.Queries) *AdminRepository {
	return &AdminRepository{q: q}
}

func (r *AdminRepository) Tx(tx ports.Transaction) ports.AdminRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &AdminRepository{q: r.q.WithTx(sqlTx)}
	}
	return r
}

func (r *AdminRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Admin, error) {
	row, err := r.q.GetAdminByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.AdminFromDB(row), nil
}

func (r *AdminRepository) FindByUsername(ctx context.Context, username string) (*domain.Admin, error) {
	row, err := r.q.GetAdminByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.AdminFromDB(row), nil
}

func (r *AdminRepository) Save(ctx context.Context, admin *domain.Admin) error {
	_, err := r.q.UpdateAdminCredentials(ctx, mappers.UpdateAdminCredentialsParamsFromDomain(admin))
	return err
}
