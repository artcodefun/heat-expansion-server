package commands

import (
	"context"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/services"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/google/uuid"
)

type DiplomacyCommands struct {
	Relationships ports.DiplomaticRelationshipRepository
	Messages      ports.DiplomaticMessageRepository
	Requests      ports.DiplomaticRequestRepository
	Operations    ports.MilitaryOperationRepository
	Users         ports.UserRepository
	UserBases     ports.UserBaseRepository
	Sectors       ports.SectorRepository
	Outbox        ports.OutboxEventRepository
	Scheduler     ports.Scheduler
	TxMgr         ports.TransactionManager
	Access        *services.AccessControlService
}

func NewDiplomacyCommands(
	relationships ports.DiplomaticRelationshipRepository,
	messages ports.DiplomaticMessageRepository,
	requests ports.DiplomaticRequestRepository,
	operations ports.MilitaryOperationRepository,
	users ports.UserRepository,
	userBases ports.UserBaseRepository,
	sectors ports.SectorRepository,
	outbox ports.OutboxEventRepository,
	scheduler ports.Scheduler,
	txMgr ports.TransactionManager,
	access *services.AccessControlService,
) *DiplomacyCommands {
	return &DiplomacyCommands{
		Relationships: relationships,
		Messages:      messages,
		Requests:      requests,
		Operations:    operations,
		Users:         users,
		UserBases:     userBases,
		Sectors:       sectors,
		Outbox:        outbox,
		Scheduler:     scheduler,
		TxMgr:         txMgr,
		Access:        access,
	}
}

func (c *DiplomacyCommands) SendInformationalMessage(ctx context.Context, actor cqrs.Actor, senderBaseID int, receiverUserID uuid.UUID, receiverBaseID *int, content domain.TranslationKey) (*uuid.UUID, error) {
	if err := c.validateDiplomaticAction(ctx, actor.UserID, senderBaseID, receiverUserID, receiverBaseID); err != nil {
		return nil, err
	}
	if !domain.IsUserSendableDiplomaticMessageContent(content) {
		return nil, cqrs.NewAppError(cqrs.KindInvalidInput, "error.application.invalid_diplomatic_message_kind")
	}

	var messageID *uuid.UUID
	err := c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		msgRepo := c.Messages.Tx(tx)
		outbox := c.Outbox.Tx(tx)

		message, err := domain.NewDiplomaticMessage(actor.UserID, receiverUserID, &senderBaseID, receiverBaseID, nil, nil, content)
		if err != nil {
			return err
		}
		if err := msgRepo.Create(ctx, message); err != nil {
			return err
		}
		if err := outbox.Save(ctx, message.PullEvents()); err != nil {
			return err
		}
		messageID = &message.ID
		return nil
	})
	if err != nil {
		return nil, err
	}
	return messageID, nil
}

func (c *DiplomacyCommands) HandleDiplomaticMessageSentEvent(ctx context.Context, ev domain.DiplomaticMessageSentEvent) error {
	if !domain.IsUserSendableDiplomaticMessageContent(ev.Content) {
		return nil
	}

	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		relRepo := c.Relationships.Tx(tx)
		rel, err := c.loadRelationship(ctx, relRepo, ev.SenderUserID, ev.ReceiverUserID)
		if err != nil {
			return err
		}
		if !rel.IsUnknown() {
			return nil
		}
		if err := rel.EstablishContact(ev.SenderUserID); err != nil {
			return err
		}

		if err := relRepo.Create(ctx, rel); err != nil {
			return err
		}
		return nil
	})
}

func (c *DiplomacyCommands) SendRequest(ctx context.Context, actor cqrs.Actor, senderBaseID int, receiverUserID uuid.UUID, receiverBaseID *int, kind domain.DiplomaticRequestKind) (*uuid.UUID, error) {
	if err := c.validateDiplomaticAction(ctx, actor.UserID, senderBaseID, receiverUserID, receiverBaseID); err != nil {
		return nil, err
	}
	if !domain.IsDiplomaticRequestKind(kind) {
		return nil, cqrs.NewAppError(cqrs.KindInvalidInput, "error.application.invalid_diplomatic_request_kind")
	}

	var requestID *uuid.UUID
	err := c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		relRepo := c.Relationships.Tx(tx)
		requestRepo := c.Requests.Tx(tx)
		outbox := c.Outbox.Tx(tx)

		rel, err := c.loadRelationship(ctx, relRepo, actor.UserID, receiverUserID)
		if err != nil {
			return err
		}

		request, err := domain.NewDiplomaticRequest(actor.UserID, receiverUserID, &senderBaseID, receiverBaseID, kind)
		if err != nil {
			return err
		}
		pendingExists, err := requestRepo.ExistsPendingByKind(ctx, actor.UserID, receiverUserID, kind)
		if err != nil {
			return repoErr(err)
		}
		if err := request.ValidateAgainstRelationship(rel, pendingExists); err != nil {
			return err
		}
		if rel.IsUnknown() {
			if err := rel.EstablishContact(actor.UserID); err != nil {
				return err
			}
			if err := relRepo.Create(ctx, rel); err != nil {
				return err
			}
		}
		if err := requestRepo.Create(ctx, request); err != nil {
			return err
		}
		if err := outbox.Save(ctx, request.PullEvents()); err != nil {
			return err
		}
		requestID = &request.ID
		return nil
	})
	if err != nil {
		return nil, err
	}
	return requestID, nil
}

func (c *DiplomacyCommands) DeclareWar(ctx context.Context, actor cqrs.Actor, senderBaseID int, receiverUserID uuid.UUID, receiverBaseID *int) (*uuid.UUID, error) {
	if err := c.validateDiplomaticAction(ctx, actor.UserID, senderBaseID, receiverUserID, receiverBaseID); err != nil {
		return nil, err
	}

	var messageID *uuid.UUID
	err := c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		relRepo := c.Relationships.Tx(tx)
		msgRepo := c.Messages.Tx(tx)
		outbox := c.Outbox.Tx(tx)

		rel, err := c.loadRelationship(ctx, relRepo, actor.UserID, receiverUserID)
		if err != nil {
			return err
		}
		if err := rel.DeclareWar(actor.UserID); err != nil {
			return err
		}
		if err := relRepo.Update(ctx, rel); err != nil {
			return err
		}

		message, err := domain.NewDiplomaticMessage(actor.UserID, receiverUserID, &senderBaseID, receiverBaseID, nil, nil, domain.DiplomaticMessageContentWarDeclaration)
		if err != nil {
			return err
		}
		if err := msgRepo.Create(ctx, message); err != nil {
			return err
		}

		events := message.PullEvents()
		if err := outbox.Save(ctx, events); err != nil {
			return err
		}
		messageID = &message.ID
		return nil
	})
	if err != nil {
		return nil, err
	}
	return messageID, nil
}

func (c *DiplomacyCommands) BreakAlliance(ctx context.Context, actor cqrs.Actor, senderBaseID int, receiverUserID uuid.UUID, receiverBaseID *int) (*uuid.UUID, error) {
	if err := c.validateDiplomaticAction(ctx, actor.UserID, senderBaseID, receiverUserID, receiverBaseID); err != nil {
		return nil, err
	}

	var messageID *uuid.UUID
	err := c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		relRepo := c.Relationships.Tx(tx)
		msgRepo := c.Messages.Tx(tx)
		outbox := c.Outbox.Tx(tx)

		rel, err := c.loadRelationship(ctx, relRepo, actor.UserID, receiverUserID)
		if err != nil {
			return err
		}
		if rel.IsUnknown() {
			return cqrs.ErrNotFound
		}
		if err := rel.BreakAlliance(actor.UserID); err != nil {
			return err
		}
		if err := relRepo.Update(ctx, rel); err != nil {
			return err
		}

		message, err := domain.NewDiplomaticMessage(actor.UserID, receiverUserID, &senderBaseID, receiverBaseID, nil, nil, domain.DiplomaticMessageContentCoalitionBreakNotice)
		if err != nil {
			return err
		}
		if err := msgRepo.Create(ctx, message); err != nil {
			return err
		}
		if err := outbox.Save(ctx, message.PullEvents()); err != nil {
			return err
		}
		messageID = &message.ID
		return nil
	})
	if err != nil {
		return nil, err
	}
	return messageID, nil
}

func (c *DiplomacyCommands) MarkChatAsRead(ctx context.Context, actor cqrs.Actor, otherUserID uuid.UUID) error {
	if actor.UserID == uuid.Nil {
		return cqrs.ErrForbidden
	}
	if err := domain.ValidateDiplomaticParticipants(actor.UserID, otherUserID); err != nil {
		return err
	}
	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		msgRepo := c.Messages.Tx(tx)
		return repoErr(msgRepo.MarkChatAsRead(ctx, actor.UserID, otherUserID))
	})
}

func (c *DiplomacyCommands) AcceptRequest(ctx context.Context, actor cqrs.Actor, senderBaseID int, requestID uuid.UUID) error {
	if err := c.Access.EnsureBaseOwnership(ctx, actor.UserID, senderBaseID); err != nil {
		return err
	}
	return c.resolveRequest(ctx, actor, senderBaseID, requestID, true)
}

func (c *DiplomacyCommands) RejectRequest(ctx context.Context, actor cqrs.Actor, senderBaseID int, requestID uuid.UUID) error {
	if err := c.Access.EnsureBaseOwnership(ctx, actor.UserID, senderBaseID); err != nil {
		return err
	}
	return c.resolveRequest(ctx, actor, senderBaseID, requestID, false)
}

func (c *DiplomacyCommands) resolveRequest(ctx context.Context, actor cqrs.Actor, senderBaseID int, requestID uuid.UUID, accepted bool) error {
	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		relRepo := c.Relationships.Tx(tx)
		msgRepo := c.Messages.Tx(tx)
		requestRepo := c.Requests.Tx(tx)
		outbox := c.Outbox.Tx(tx)

		request, err := requestRepo.FindByIDForUpdate(ctx, requestID)
		if err != nil {
			return repoErr(err)
		}

		rel, err := c.loadRelationship(ctx, relRepo, request.SenderUserID, request.ReceiverUserID)
		if err != nil {
			return err
		}
		if rel.IsUnknown() {
			return cqrs.ErrNotFound
		}

		var responseContent domain.TranslationKey
		if accepted {
			responseContent, err = request.Accept(actor.UserID, rel)
		} else {
			responseContent, err = request.Reject(actor.UserID)
		}
		if err != nil {
			return err
		}

		if err := requestRepo.Update(ctx, request); err != nil {
			return err
		}
		if accepted {
			if err := relRepo.Update(ctx, rel); err != nil {
				return err
			}
		}
		replyToMessageID, err := c.findRequestMessageID(ctx, msgRepo, request)
		if err != nil {
			return err
		}
		response, err := domain.NewDiplomaticMessage(actor.UserID, request.SenderUserID, &senderBaseID, request.SenderBaseID, &request.ID, replyToMessageID, responseContent)
		if err != nil {
			return err
		}
		if err := msgRepo.Create(ctx, response); err != nil {
			return err
		}

		var events []domain.DomainEvent
		events = append(events, request.PullEvents()...)
		events = append(events, response.PullEvents()...)
		return outbox.Save(ctx, events)
	})
}

func (c *DiplomacyCommands) HandleDiplomaticRequestCreatedEvent(ctx context.Context, ev domain.DiplomaticRequestCreatedEvent) error {
	var expiresAt int64
	err := c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		msgRepo := c.Messages.Tx(tx)
		requestRepo := c.Requests.Tx(tx)
		outbox := c.Outbox.Tx(tx)

		request, err := requestRepo.FindByID(ctx, ev.RequestID)
		if err != nil {
			return repoErr(err)
		}
		expiresAt = request.ExpiresAt
		messageContent := request.MessageContent()
		if messageContent == "" {
			return cqrs.NewAppError(cqrs.KindInvalidInput, "error.application.invalid_diplomatic_request_kind")
		}
		exists, err := msgRepo.ExistsByRequestAndContent(ctx, request.ID, messageContent)
		if err != nil {
			return err
		}
		if exists {
			return nil
		}

		message, err := domain.NewDiplomaticMessage(request.SenderUserID, request.ReceiverUserID, request.SenderBaseID, request.ReceiverBaseID, &request.ID, nil, messageContent)
		if err != nil {
			return err
		}
		if err := msgRepo.Create(ctx, message); err != nil {
			return err
		}
		if err := outbox.Save(ctx, message.PullEvents()); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	if expiresAt == 0 {
		return nil
	}
	return c.Scheduler.Schedule(ctx, ports.ExpireDiplomaticRequestJob{RequestID: ev.RequestID}, expiresAt)
}

func (c *DiplomacyCommands) HandleMilitaryOperationResolvedEvent(ctx context.Context, ev domain.MilitaryOperationResolvedEvent) error {
	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		opRepo := c.Operations.Tx(tx)
		relRepo := c.Relationships.Tx(tx)
		baseRepo := c.UserBases.Tx(tx)
		sectorRepo := c.Sectors.Tx(tx)

		op, err := opRepo.FindByID(ctx, ev.OperationID)
		if err != nil {
			return repoErr(err)
		}
		if op.Type != domain.MilitaryOperationTypeAttack {
			return nil
		}
		occType, err := sectorRepo.GetLocationTypeByCoordinates(ctx, op.TargetCoordinates.X, op.TargetCoordinates.Y)
		if err != nil {
			return repoErr(err)
		}
		if occType != domain.LocationTypeUserBase {
			return nil
		}
		defenderBase, err := baseRepo.FindByCoordinates(ctx, op.TargetCoordinates.X, op.TargetCoordinates.Y)
		if err != nil {
			return repoErr(err)
		}
		if op.OwnerUserID == defenderBase.UserID {
			return nil
		}

		rel, err := c.loadRelationship(ctx, relRepo, op.OwnerUserID, defenderBase.UserID)
		if err != nil {
			return err
		}
		if !rel.IsUnknown() {
			return nil
		}
		if err := rel.EscalateToWar(op.OwnerUserID); err != nil {
			return err
		}
		return relRepo.Create(ctx, rel)
	})
}

func (c *DiplomacyCommands) HandleExpireDiplomaticRequestJob(ctx context.Context, cmd ports.ExpireDiplomaticRequestJob) error {
	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		requestRepo := c.Requests.Tx(tx)
		request, err := requestRepo.FindByIDForUpdate(ctx, cmd.RequestID)
		if err != nil {
			if errors.Is(err, ports.ErrNotFound) {
				return nil
			}
			return repoErr(err)
		}
		if !request.CanExpire() {
			return nil
		}
		request.Expire()
		return requestRepo.Update(ctx, request)
	})
}

func (c *DiplomacyCommands) loadRelationship(ctx context.Context, repo ports.DiplomaticRelationshipRepository, senderUserID, receiverUserID uuid.UUID) (*domain.DiplomaticRelationship, error) {
	rel, err := repo.FindBetweenUsersForUpdate(ctx, senderUserID, receiverUserID)
	if err == nil {
		return rel, nil
	}
	if !errors.Is(err, ports.ErrNotFound) {
		return nil, repoErr(err)
	}
	return domain.NewUnknownRelationship(senderUserID, receiverUserID)
}

func (c *DiplomacyCommands) findRequestMessageID(ctx context.Context, repo ports.DiplomaticMessageRepository, request *domain.DiplomaticRequest) (*uuid.UUID, error) {
	message, err := repo.FindByRequestAndContent(ctx, request.ID, request.MessageContent())
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &message.ID, nil
}

func (c *DiplomacyCommands) validateDiplomaticAction(ctx context.Context, actorUserID uuid.UUID, senderBaseID int, receiverUserID uuid.UUID, receiverBaseID *int) error {
	if err := c.Access.EnsureBaseOwnership(ctx, actorUserID, senderBaseID); err != nil {
		return err
	}
	if err := domain.ValidateDiplomaticParticipants(actorUserID, receiverUserID); err != nil {
		return err
	}
	return c.validateReceiverTarget(ctx, receiverUserID, receiverBaseID)
}

func (c *DiplomacyCommands) validateReceiverTarget(ctx context.Context, receiverUserID uuid.UUID, receiverBaseID *int) error {
	if _, err := c.Users.FindByID(ctx, receiverUserID); err != nil {
		return repoErr(err)
	}
	if receiverBaseID == nil {
		return nil
	}
	ownerID, err := c.UserBases.GetOwnerID(ctx, *receiverBaseID)
	if err != nil {
		return repoErr(err)
	}
	if ownerID != receiverUserID {
		return cqrs.NewAppError(cqrs.KindInvalidInput, "error.application.invalid_diplomatic_receiver_target")
	}
	return nil
}
