package bootstrap

import (
	"context"
	"database/sql"
	"time"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	dbgen "github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/repo"
)

// seedPeriodicJobs ensures recurring jobs exist once at startup.
func seedPeriodicJobs(ctx context.Context, db *sql.DB) error {
	q := dbgen.New(db)
	schedulerRepo := repo.NewScheduledJobRepo(q)

	return seedBlackMarketRefreshJob(ctx, schedulerRepo)
}

func seedBlackMarketRefreshJob(ctx context.Context, scheduledJobs *repo.ScheduledJobRepo) error {
	job := ports.RefreshBlackMarketOffersJob{}
	now := time.Now().Unix()
	inserted, err := scheduledJobs.InsertIfNotExists(ctx, job, now+60, now)
	if err != nil {
		return err
	}
	if !inserted {
		return nil
	}
	return nil
}
