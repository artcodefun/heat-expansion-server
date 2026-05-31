package commands

import (
	"context"
	"math/rand"
	"time"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/services"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
)

const (
	blackMarketLimitedOfferMaxActivePerKind = 2
	blackMarketLimitedOfferMaxActiveTotal   = 2
	blackMarketRefreshPeriod                = 12 * time.Hour
)

type BlackMarketCommands struct {
	UserRepo         ports.UserRepository
	BaseRepo         ports.UserBaseRepository
	OfferRepo        ports.BlackMarketOfferRepository
	BuildPrototypes  ports.BuildPrototypeRepository
	ArmyPrototypes   ports.ArmyPrototypeRepository
	StoragePrototype ports.StoragePrototypeRepository
	Outbox           ports.OutboxEventRepository
	Scheduler        ports.Scheduler
	TxMgr            ports.TransactionManager
	Access           *services.AccessControlService
	BlackMarket      *domain.BlackMarketService
}

func NewBlackMarketCommands(
	userRepo ports.UserRepository,
	baseRepo ports.UserBaseRepository,
	offerRepo ports.BlackMarketOfferRepository,
	buildPrototypes ports.BuildPrototypeRepository,
	armyPrototypes ports.ArmyPrototypeRepository,
	storagePrototypes ports.StoragePrototypeRepository,
	outbox ports.OutboxEventRepository,
	scheduler ports.Scheduler,
	txMgr ports.TransactionManager,
	access *services.AccessControlService,
) *BlackMarketCommands {
	return &BlackMarketCommands{
		UserRepo:         userRepo,
		BaseRepo:         baseRepo,
		OfferRepo:        offerRepo,
		BuildPrototypes:  buildPrototypes,
		ArmyPrototypes:   armyPrototypes,
		StoragePrototype: storagePrototypes,
		Outbox:           outbox,
		Scheduler:        scheduler,
		TxMgr:            txMgr,
		Access:           access,
		BlackMarket:      domain.NewBlackMarketService(),
	}
}

func (c *BlackMarketCommands) PurchaseResources(ctx context.Context, actor cqrs.Actor, baseID int, resourceType domain.ResourceType, crystals int) error {
	if err := c.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return err
	}

	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		uRepo := c.UserRepo.Tx(tx)
		bRepo := c.BaseRepo.Tx(tx)

		user, err := uRepo.FindByIDForUpdate(ctx, actor.UserID)
		if err != nil {
			return repoErr(err)
		}
		base, err := bRepo.FindByIDForUpdate(ctx, baseID)
		if err != nil {
			return repoErr(err)
		}

		if err := c.BlackMarket.PurchaseResources(user, base, resourceType, crystals); err != nil {
			return err
		}
		if err := uRepo.Update(ctx, user); err != nil {
			return repoErr(err)
		}
		if err := bRepo.Update(ctx, base); err != nil {
			return repoErr(err)
		}
		if err := c.Outbox.Tx(tx).Save(ctx, user.EventProducer.PullEvents()); err != nil {
			return repoErr(err)
		}
		return nil
	})
}

func (c *BlackMarketCommands) PurchaseOffer(ctx context.Context, actor cqrs.Actor, baseID int, offerID int64, quantity int) error {
	if err := c.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return err
	}

	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		uRepo := c.UserRepo.Tx(tx)
		bRepo := c.BaseRepo.Tx(tx)
		oRepo := c.OfferRepo.Tx(tx)

		user, err := uRepo.FindByIDForUpdate(ctx, actor.UserID)
		if err != nil {
			return repoErr(err)
		}
		base, err := bRepo.FindByIDForUpdate(ctx, baseID)
		if err != nil {
			return repoErr(err)
		}
		offer, err := oRepo.FindByIDForUpdate(ctx, offerID)
		if err != nil {
			return repoErr(err)
		}

		switch offer.Kind {
		case domain.BlackMarketOfferKindBuilding:
			proto, err := c.BuildPrototypes.FindPrototypeByID(ctx, offer.PrototypeID)
			if err != nil {
				return repoErr(err)
			}
			if err := c.BlackMarket.PurchaseBuildingOffer(user, base, *offer, proto); err != nil {
				return err
			}
		case domain.BlackMarketOfferKindArmy:
			proto, err := c.ArmyPrototypes.FindPrototypeByID(ctx, offer.PrototypeID)
			if err != nil {
				return repoErr(err)
			}
			if err := c.BlackMarket.PurchaseArmyOffer(user, base, *offer, proto, quantity); err != nil {
				return err
			}
		case domain.BlackMarketOfferKindStorage:
			proto, err := c.StoragePrototype.FindPrototypeByID(ctx, offer.PrototypeID)
			if err != nil {
				return repoErr(err)
			}
			if err := c.BlackMarket.PurchaseStorageOffer(user, base, *offer, proto); err != nil {
				return err
			}
		default:
			return cqrs.NewAppErrorWithParams(cqrs.KindInvalidInput, "error.application.black_market.offer_kind_invalid", map[string]any{"kind": offer.Kind})
		}

		if err := uRepo.Update(ctx, user); err != nil {
			return repoErr(err)
		}
		if err := bRepo.Update(ctx, base); err != nil {
			return repoErr(err)
		}
		if err := c.Outbox.Tx(tx).Save(ctx, user.EventProducer.PullEvents()); err != nil {
			return repoErr(err)
		}
		return nil
	})
}

func (c *BlackMarketCommands) HandleRefreshBlackMarketOffersJob(ctx context.Context, job ports.RefreshBlackMarketOffersJob) error {
	now := time.Now().Unix()
	err := c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		oRepo := c.OfferRepo.Tx(tx)

		activeOffers, err := oRepo.ListActiveLimitedOffers(ctx, now)
		if err != nil {
			return err
		}

		activeCountByKind := make(map[domain.BlackMarketOfferKind]int)
		for _, offer := range activeOffers {
			activeCountByKind[offer.Kind]++
		}
		activeCountTotal := len(activeOffers)
		promotedKindsThisRun := make(map[domain.BlackMarketOfferKind]struct{})

		expiredOffers, err := oRepo.ListExpiredLimitedOffers(ctx, now)
		if err != nil {
			return err
		}

		for _, offer := range expiredOffers {
			if activeCountTotal >= blackMarketLimitedOfferMaxActiveTotal {
				break
			}
			if activeCountByKind[offer.Kind] >= blackMarketLimitedOfferMaxActivePerKind {
				continue
			}
			if _, alreadyPromoted := promotedKindsThisRun[offer.Kind]; alreadyPromoted {
				continue
			}
			expiresAt := now + domain.BlackMarketLimitedOfferDurationSeconds
			offer.EndsAt = &expiresAt
			if err := oRepo.Update(ctx, offer); err != nil {
				return err
			}
			activeCountByKind[offer.Kind]++
			activeCountTotal++
			promotedKindsThisRun[offer.Kind] = struct{}{}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return c.rescheduleRefreshJob(ctx)
}

func (c *BlackMarketCommands) rescheduleRefreshJob(ctx context.Context) error {
	jitter := int64(rand.Intn(300))
	return c.Scheduler.Schedule(ctx, ports.RefreshBlackMarketOffersJob{}, time.Now().Unix()+int64(blackMarketRefreshPeriod.Seconds())+jitter)
}
