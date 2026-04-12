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

type DiplomaticMessageRepo struct {
	q *gen.Queries
}

func NewDiplomaticMessageRepo(q *gen.Queries) *DiplomaticMessageRepo {
	return &DiplomaticMessageRepo{q: q}
}

func (r *DiplomaticMessageRepo) Tx(tx ports.Transaction) ports.DiplomaticMessageRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &DiplomaticMessageRepo{q: r.q.WithTx(sqlTx)}
	}
	return r
}

func (r *DiplomaticMessageRepo) Create(ctx context.Context, message *domain.DiplomaticMessage) error {
	return r.q.InsertDiplomaticMessage(ctx, mappers.InsertDiplomaticMessageParamsFromDomain(message))
}

func (r *DiplomaticMessageRepo) ExistsByRequestAndContent(ctx context.Context, requestID uuid.UUID, content domain.TranslationKey) (bool, error) {
	return r.q.ExistsDiplomaticMessageByRequestAndContent(ctx, gen.ExistsDiplomaticMessageByRequestAndContentParams{
		RequestID: uuid.NullUUID{UUID: requestID, Valid: true},
		Content:   content,
	})
}

func (r *DiplomaticMessageRepo) FindByRequestAndContent(ctx context.Context, requestID uuid.UUID, content domain.TranslationKey) (*domain.DiplomaticMessage, error) {
	row, err := r.q.GetDiplomaticMessageByRequestAndContent(ctx, gen.GetDiplomaticMessageByRequestAndContentParams{
		RequestID: uuid.NullUUID{UUID: requestID, Valid: true},
		Content:   content,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.DiplomaticMessageFromDB(row), nil
}

func (r *DiplomaticMessageRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.DiplomaticMessage, error) {
	row, err := r.q.GetDiplomaticMessage(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.DiplomaticMessageFromDB(row), nil
}

func (r *DiplomaticMessageRepo) MarkChatAsRead(ctx context.Context, receiverUserID, senderUserID uuid.UUID) error {
	return r.q.MarkDiplomaticChatAsRead(ctx, gen.MarkDiplomaticChatAsReadParams{
		ReceiverUserID: receiverUserID,
		SenderUserID:   senderUserID,
	})
}
