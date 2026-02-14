package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/mappers"
	"github.com/google/uuid"
)

type UserRepo struct {
	q *gen.Queries
}

func NewUserRepo(q *gen.Queries) *UserRepo { return &UserRepo{q: q} }

func (r *UserRepo) Tx(tx ports.Transaction) ports.UserRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &UserRepo{q: r.q.WithTx(sqlTx)}
	}
	return r
}

func (r *UserRepo) Create(user *domain.User) error {
	err := r.q.InsertUser(context.Background(), gen.InsertUserParams{
		ID:       user.ID,
		Name:     user.Name,
		Crystals: int32(user.Crystals),
	})
	return err
}

func (r *UserRepo) FindByID(id uuid.UUID) (*domain.User, error) {
	u, err := r.q.GetUserByID(context.Background(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.UserFromDB(u), nil
}

// FindByIDForUpdate uses a FOR UPDATE lock. Requires a transaction-bound repo.
func (r *UserRepo) FindByIDForUpdate(id uuid.UUID) (*domain.User, error) {
	// sqlc does not generate a FOR UPDATE variant yet; placeholder for future query.
	return r.FindByID(id)
}

func (r *UserRepo) Update(user *domain.User) error {
	err := r.q.UpdateUser(context.Background(), gen.UpdateUserParams{
		ID:       user.ID,
		Name:     user.Name,
		Crystals: int32(user.Crystals),
	})
	return err
}

func (r *UserRepo) Delete(id uuid.UUID) error {
	// No delete query defined yet; placeholder error to avoid silent omission.
	return errors.New("Delete not implemented for users")
}
