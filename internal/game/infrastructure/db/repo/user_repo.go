package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/mappers"
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
	id, err := r.q.InsertUser(context.Background(), gen.InsertUserParams{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Crystals:     int32(user.Crystals),
	})
	if err != nil {
		return err
	}
	user.ID = int(id)
	return nil
}

func (r *UserRepo) FindByID(id int) (*domain.User, error) {
	u, err := r.q.GetUserByID(context.Background(), int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.UserFromDB(u), nil
}

// FindByIDForUpdate uses a FOR UPDATE lock. Requires a transaction-bound repo.
func (r *UserRepo) FindByIDForUpdate(id int) (*domain.User, error) {
	// sqlc does not generate a FOR UPDATE variant yet; placeholder for future query.
	return r.FindByID(id)
}

func (r *UserRepo) FindByEmail(email string) (*domain.User, error) {
	u, err := r.q.GetUserByEmail(context.Background(), email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.UserFromDB(u), nil
}

func (r *UserRepo) Update(user *domain.User) error {
	err := r.q.UpdateUser(context.Background(), gen.UpdateUserParams{
		ID:           int64(user.ID),
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Crystals:     int32(user.Crystals),
	})
	return err
}

func (r *UserRepo) Delete(id int) error {
	// No delete query defined yet; placeholder error to avoid silent omission.
	return errors.New("Delete not implemented for users")
}
