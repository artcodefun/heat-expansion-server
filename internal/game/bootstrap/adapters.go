package bootstrap

import (
	"context"
	"database/sql"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	contentgen "github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/content"
	dbgen "github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/gen"
	repo "github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/repo"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/i18n"
	jobs "github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/jobs"
	readgen "github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/gen"
	readrepo "github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/repo"
	platformevents "github.com/artcodefun/heat-expansion-server/internal/platform/events"
	"github.com/artcodefun/heat-expansion-server/internal/platform/security"
)

// Adapters wires secondary adapters (repositories, tx manager) implementing core ports.
type Adapters struct {
	// Repositories
	Users                   ports.UserRepository
	UserBases               ports.UserBaseRepository
	Sectors                 ports.SectorRepository
	ResourceLocations       ports.ResourceLocationRepository
	DangerousLocations      ports.DangerousLocationRepository
	ArmyPrototypes          ports.ArmyPrototypeRepository
	BuildPrototypes         ports.BuildPrototypeRepository
	StoragePrototypes       ports.StoragePrototypeRepository
	TechPrototypes          ports.TechPrototypeRepository
	BlackMarketOffers       ports.BlackMarketOfferRepository
	MilitaryOps             ports.MilitaryOperationRepository
	TradeOps                ports.TradeOperationRepository
	ScanReports             ports.ScanReportRepository
	Activities              ports.ActivityRepository
	RadarThreats            ports.RadarThreatRepository
	OutboxEvents            ports.OutboxEventRepository
	Alerts                  ports.AlertRepository
	CrystalCredits          ports.CrystalCreditsRepository
	DiplomaticRelationships ports.DiplomaticRelationshipRepository
	DiplomaticMessages      ports.DiplomaticMessageRepository
	DiplomaticRequests      ports.DiplomaticRequestRepository

	// Read Repositories (read-store / projections)
	BaseRead           ports.BaseReadRepository
	BuildingRead       ports.BuildingReadRepository
	ArmyRead           ports.ArmyReadRepository
	ArmyPrototypeRead  ports.ArmyPrototypeReadRepository
	BuildPrototypeRead ports.BuildPrototypeReadRepository
	StorageRead        ports.StorageReadRepository
	TechRead           ports.TechReadRepository
	OperationRead      ports.OperationReadRepository
	TradeOperationRead ports.TradeOperationReadRepository
	BlackMarketRead    ports.BlackMarketReadRepository
	ActivityRead       ports.ActivityReadRepository
	SectorRead         ports.SectorReadRepository
	RadarRead          ports.RadarReadRepository
	UserRead           ports.UserReadRepository
	AlertRead          ports.AlertReadRepository
	DiplomacyRead      ports.DiplomacyReadRepository

	// Infra
	TxMgr      ports.TransactionManager
	Tokens     ports.TokenValidator
	Events     ports.EventPublisher
	Scheduler  ports.Scheduler
	Content    ports.ContentGenerator
	Translator ports.Translator

	// translator keeps the concrete type so Setup can run its startup I/O.
	translator *i18n.SimpleTranslator
}

// Setup performs the adapters' one-time startup I/O (loading content
// translations from the database). Unlike Start(ctx) on the broker adapters it
// is not long-running: it returns once the load completes. It must complete
// before the HTTP server starts serving requests, so handlers never observe a
// partially loaded translator.
func (a *Adapters) Setup(ctx context.Context) error {
	if err := a.translator.LoadFromRepo(ctx); err != nil {
		return err
	}
	return nil
}

func NewAdapters(db *sql.DB, staticBaseURL string, jwtPublicKeyPEM string) (*Adapters, error) {
	q := dbgen.New(db)
	rq := readgen.New(db)

	publisher := platformevents.NewSimplePublisher[domain.DomainEvent]()
	txMgr := repo.NewDBTxManager(db)
	schedulerRepo := repo.NewScheduledJobRepo(q)
	scheduler := jobs.NewDBScheduler(txMgr, schedulerRepo)

	generator := contentgen.NewSimpleGenerator(staticBaseURL)
	tokens, err := security.NewSimpleTokenValidator(jwtPublicKeyPEM)
	if err != nil {
		return nil, err
	}

	translationRepo := repo.NewTranslationRepo(q)
	translator, err := i18n.NewSimpleTranslator(translationRepo)
	if err != nil {
		return nil, err
	}

	baseRead := readrepo.NewBaseReadRepo(rq)
	sectorRead := readrepo.NewSectorReadRepo(rq, baseRead)
	opRead := readrepo.NewOperationReadRepo(rq, sectorRead)
	tradeOpRead := readrepo.NewTradeOperationReadRepo(rq, sectorRead)
	radarRead := readrepo.NewRadarThreatReadRepo(rq)

	armyProtoRepo := repo.NewArmyPrototypeRepo(q)
	buildProtoRepo := repo.NewBuildPrototypeRepo(q)

	return &Adapters{
		// Repositories
		Users:                   repo.NewUserRepo(q),
		UserBases:               repo.NewUserBaseRepo(q),
		Sectors:                 repo.NewSectorRepo(q),
		ResourceLocations:       repo.NewResourceLocationRepo(q, armyProtoRepo, buildProtoRepo),
		DangerousLocations:      repo.NewDangerousLocationRepo(q, armyProtoRepo, buildProtoRepo),
		ArmyPrototypes:          armyProtoRepo,
		BuildPrototypes:         buildProtoRepo,
		StoragePrototypes:       repo.NewStoragePrototypeRepo(q),
		TechPrototypes:          repo.NewTechPrototypeRepo(q),
		BlackMarketOffers:       repo.NewBlackMarketOfferRepo(q),
		MilitaryOps:             repo.NewMilitaryOperationRepo(q),
		TradeOps:                repo.NewTradeOperationRepo(q),
		ScanReports:             repo.NewScanReportRepo(q),
		Activities:              repo.NewActivityRepo(q),
		RadarThreats:            repo.NewRadarThreatRepo(q),
		OutboxEvents:            repo.NewOutboxEventRepo(q),
		Alerts:                  repo.NewAlertRepo(q),
		CrystalCredits:          repo.NewCrystalCreditsRepo(q),
		DiplomaticRelationships: repo.NewDiplomaticRelationshipRepo(q),
		DiplomaticMessages:      repo.NewDiplomaticMessageRepo(q),
		DiplomaticRequests:      repo.NewDiplomaticRequestRepo(q),

		// Read side
		BaseRead:           baseRead,
		BuildingRead:       readrepo.NewBuildReadRepo(rq),
		ArmyRead:           readrepo.NewArmyReadRepo(rq),
		ArmyPrototypeRead:  readrepo.NewPrototypeReadRepo(rq),
		BuildPrototypeRead: readrepo.NewPrototypeReadRepo(rq),
		StorageRead:        readrepo.NewStorageReadRepo(rq),
		TechRead:           readrepo.NewTechReadRepo(rq),
		OperationRead:      opRead,
		TradeOperationRead: tradeOpRead,
		BlackMarketRead:    readrepo.NewBlackMarketReadRepo(rq),
		ActivityRead:       readrepo.NewActivityReadRepo(rq, opRead, tradeOpRead, sectorRead, radarRead),
		SectorRead:         sectorRead,
		RadarRead:          radarRead,
		UserRead:           readrepo.NewUserReadRepo(rq),
		AlertRead:          readrepo.NewAlertReadRepository(rq),
		DiplomacyRead:      readrepo.NewDiplomacyReadRepo(rq, baseRead),

		// Infra
		TxMgr:      txMgr,
		Tokens:     tokens,
		Events:     publisher,
		Scheduler:  scheduler,
		Content:    generator,
		Translator: translator,
		translator: translator,
	}, nil
}
