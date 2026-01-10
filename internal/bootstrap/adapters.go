package bootstrap

import (
	"database/sql"

	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	contentgen "github.com/artcodefun/heat-expansion-api/internal/infrastructure/content"
	dbgen "github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/gen"
	repo "github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/repo"
	events "github.com/artcodefun/heat-expansion-api/internal/infrastructure/events"
	jobs "github.com/artcodefun/heat-expansion-api/internal/infrastructure/jobs"
	readgen "github.com/artcodefun/heat-expansion-api/internal/infrastructure/readstore/gen"
	readrepo "github.com/artcodefun/heat-expansion-api/internal/infrastructure/readstore/repo"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/security"
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
	OutboxEvents       ports.OutboxEventRepository

	// Read Repositories (read-store / projections)
	BaseRead      ports.BaseReadRepository
	BuildingRead  ports.BuildingReadRepository
	ArmyRead      ports.ArmyReadRepository
	StorageRead   ports.StorageReadRepository
	TechRead      ports.TechReadRepository
	OperationRead ports.OperationReadRepository
	ActivityRead  ports.ActivityReadRepository
	SectorRead    ports.SectorReadRepository
	UserRead      ports.UserReadRepository

	// Infra
	TxMgr     ports.TransactionManager
	Events    ports.EventPublisher
	Scheduler ports.Scheduler
	Hasher    ports.PasswordHasher
	Tokens    ports.TokenProvider
	Content   ports.ContentGenerator
}

func NewAdapters(db *sql.DB, jwtSecret, contentDir, staticBaseURL string) (*Adapters, error) {
	q := dbgen.New(db)
	rq := readgen.New(db)

	// Simple in-process publisher + DB-backed scheduler (durable jobs).
	publisher := events.NewSimplePublisher()
	txMgr := repo.NewDBTxManager(db)
	schedulerRepo := repo.NewScheduledJobRepo(q)
	scheduler := jobs.NewDBScheduler(txMgr, schedulerRepo)
	// Security + content adapters (dev-friendly defaults)
	hasher := security.NewSimpleHasher()
	tokens := security.NewSimpleTokenProvider(jwtSecret)
	generator := contentgen.NewSimpleGenerator(contentDir, staticBaseURL)

	opRead := readrepo.NewOperationReadRepo(rq)
	sectorRead := readrepo.NewSectorReadRepo(rq)
	activityRead := readrepo.NewActivityReadRepo(rq, opRead, sectorRead)

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
		OutboxEvents:       repo.NewOutboxEventRepo(q),
		// Read side
		BaseRead:      readrepo.NewBaseReadRepo(rq),
		BuildingRead:  readrepo.NewBuildReadRepo(rq),
		ArmyRead:      readrepo.NewArmyReadRepo(rq),
		StorageRead:   readrepo.NewStorageReadRepo(rq),
		TechRead:      readrepo.NewTechReadRepo(rq),
		OperationRead: opRead,
		ActivityRead:  activityRead,
		SectorRead:    sectorRead,
		UserRead:      readrepo.NewUserReadRepo(rq),
		TxMgr:         txMgr,
		Events:        publisher,
		Scheduler:     scheduler,
		Hasher:        hasher,
		Tokens:        tokens,
		Content:       generator,
	}, nil
}
