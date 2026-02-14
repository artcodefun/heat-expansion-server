package bootstrap

import (
	"database/sql"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	contentgen "github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/content"
	dbgen "github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/gen"
	repo "github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/repo"
	events "github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/events"
	jobs "github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/jobs"
	readgen "github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/gen"
	readrepo "github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/repo"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/security"
)

// Adapters wires secondary adapters (repositories, tx manager) implementing core ports.
type Adapters struct {
	// Repositories
	Users              ports.UserRepository
	UserBases          ports.UserBaseRepository
	Sectors            ports.SectorRepository
	ResourceLocations  ports.ResourceLocationRepository
	DangerousLocations ports.DangerousLocationRepository
	ArmyPrototypes     ports.ArmyPrototypeRepository
	BuildPrototypes    ports.BuildPrototypeRepository
	StoragePrototypes  ports.StoragePrototypeRepository
	TechPrototypes     ports.TechPrototypeRepository
	MilitaryOps        ports.MilitaryOperationRepository
	ScanReports        ports.ScanReportRepository
	Activities         ports.ActivityRepository
	RadarThreats       ports.RadarThreatRepository
	OutboxEvents       ports.OutboxEventRepository
	Alerts             ports.AlertRepository

	// Read Repositories (read-store / projections)
	BaseRead      ports.BaseReadRepository
	BuildingRead  ports.BuildingReadRepository
	ArmyRead      ports.ArmyReadRepository
	StorageRead   ports.StorageReadRepository
	TechRead      ports.TechReadRepository
	OperationRead ports.OperationReadRepository
	ActivityRead  ports.ActivityReadRepository
	SectorRead    ports.SectorReadRepository
	RadarRead     ports.RadarReadRepository
	UserRead      ports.UserReadRepository
	AlertRead     ports.AlertReadRepository

	// Infra
	TxMgr     ports.TransactionManager
	Tokens    ports.TokenValidator
	Events    ports.EventPublisher
	Scheduler ports.Scheduler
	Content   ports.ContentGenerator
}

func NewAdapters(db *sql.DB, staticBaseURL string, jwtSecret string) (*Adapters, error) {
	q := dbgen.New(db)
	rq := readgen.New(db)

	publisher := events.NewSimplePublisher()
	txMgr := repo.NewDBTxManager(db)
	schedulerRepo := repo.NewScheduledJobRepo(q)
	scheduler := jobs.NewDBScheduler(txMgr, schedulerRepo)

	generator := contentgen.NewSimpleGenerator(staticBaseURL)
	tokens := security.NewSimpleTokenValidator(jwtSecret)

	sectorRead := readrepo.NewSectorReadRepo(rq)
	opRead := readrepo.NewOperationReadRepo(rq, sectorRead)
	radarRead := readrepo.NewRadarThreatReadRepo(rq)
	activityRead := readrepo.NewActivityReadRepo(rq, opRead, sectorRead, radarRead)

	armyProtoRepo := repo.NewArmyPrototypeRepo(q)
	buildProtoRepo := repo.NewBuildPrototypeRepo(q)

	return &Adapters{
		Users:              repo.NewUserRepo(q),
		UserBases:          repo.NewUserBaseRepo(q),
		Sectors:            repo.NewSectorRepo(q),
		ResourceLocations:  repo.NewResourceLocationRepo(q, armyProtoRepo, buildProtoRepo),
		DangerousLocations: repo.NewDangerousLocationRepo(q, armyProtoRepo, buildProtoRepo),
		ArmyPrototypes:     armyProtoRepo,
		BuildPrototypes:    buildProtoRepo,
		StoragePrototypes:  repo.NewStoragePrototypeRepo(q),
		TechPrototypes:     repo.NewTechPrototypeRepo(q),
		MilitaryOps:        repo.NewMilitaryOperationRepo(q),
		ScanReports:        repo.NewScanReportRepo(q),
		Activities:         repo.NewActivityRepo(q),
		RadarThreats:       repo.NewRadarThreatRepo(q),
		OutboxEvents:       repo.NewOutboxEventRepo(q),
		Alerts:             repo.NewAlertRepo(q),
		// Read side
		BaseRead:      readrepo.NewBaseReadRepo(rq),
		BuildingRead:  readrepo.NewBuildReadRepo(rq),
		ArmyRead:      readrepo.NewArmyReadRepo(rq),
		StorageRead:   readrepo.NewStorageReadRepo(rq),
		TechRead:      readrepo.NewTechReadRepo(rq),
		OperationRead: opRead,
		ActivityRead:  activityRead,
		SectorRead:    sectorRead,
		RadarRead:     radarRead,
		UserRead:      readrepo.NewUserReadRepo(rq),
		AlertRead:     readrepo.NewAlertReadRepository(rq),
		TxMgr:         txMgr,
		Tokens:        tokens,
		Events:        publisher,
		Scheduler:     scheduler,
		Content:       generator,
	}, nil
}
