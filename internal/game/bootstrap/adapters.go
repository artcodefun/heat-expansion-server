package bootstrap

import (
	"database/sql"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	contentgen "github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/content"
	dbgen "github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/gen"
	repo "github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/repo"
	events "github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/events"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/i18n"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/i18n/locales"
	jobs "github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/jobs"
	readgen "github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/gen"
	readrepo "github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/repo"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/security"
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
	DiplomaticRelationships ports.DiplomaticRelationshipRepository
	DiplomaticMessages      ports.DiplomaticMessageRepository
	DiplomaticRequests      ports.DiplomaticRequestRepository

	// Read Repositories (read-store / projections)
	BaseRead           ports.BaseReadRepository
	BuildingRead       ports.BuildingReadRepository
	ArmyRead           ports.ArmyReadRepository
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
}

func NewAdapters(db *sql.DB, staticBaseURL string, jwtSecret string, i18nPath string) (*Adapters, error) {
	q := dbgen.New(db)
	rq := readgen.New(db)

	publisher := events.NewSimplePublisher()
	txMgr := repo.NewDBTxManager(db)
	schedulerRepo := repo.NewScheduledJobRepo(q)
	scheduler := jobs.NewDBScheduler(txMgr, schedulerRepo)

	generator := contentgen.NewSimpleGenerator(staticBaseURL)
	tokens := security.NewSimpleTokenValidator(jwtSecret)

	translator := i18n.NewJSONTranslator()
	// 1. Load systemic locales (Embedded in binary)
	if err := translator.LoadFromFS(locales.Files, "."); err != nil {
		return nil, err
	}

	// 2. Load content locales (External directory provided via bootstrap)
	if i18nPath != "" {
		if err := translator.LoadFromDir(i18nPath); err != nil {
			return nil, err
		}
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
		DiplomaticRelationships: repo.NewDiplomaticRelationshipRepo(q),
		DiplomaticMessages:      repo.NewDiplomaticMessageRepo(q),
		DiplomaticRequests:      repo.NewDiplomaticRequestRepo(q),

		// Read side
		BaseRead:           baseRead,
		BuildingRead:       readrepo.NewBuildReadRepo(rq),
		ArmyRead:           readrepo.NewArmyReadRepo(rq),
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
	}, nil
}
