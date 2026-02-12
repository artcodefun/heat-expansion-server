package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/mappers"
)

type ScanReportRepo struct {
	q *gen.Queries
}

func NewScanReportRepo(q *gen.Queries) *ScanReportRepo { return &ScanReportRepo{q: q} }

func (r *ScanReportRepo) Tx(tx ports.Transaction) ports.ScanReportRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &ScanReportRepo{q: r.q.WithTx(sqlTx)}
	}
	return r
}

func (r *ScanReportRepo) Create(report *domain.SectorScanReport) error {
	id, err := r.q.InsertScanReport(context.Background(), mappers.InsertScanReportParamsFromDomain(report))
	if err != nil {
		return err
	}
	report.ID = int(id)
	return nil
}

func (r *ScanReportRepo) FindByID(id int) (*domain.SectorScanReport, error) {
	row, err := r.q.GetScanReportByID(context.Background(), int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.ScanReportFromDB(row), nil
}

func (r *ScanReportRepo) RecentReportExistsByScanner(scannerID uuid.UUID, since int64) (bool, error) {
	exists, err := r.q.RecentReportExistsByScanner(context.Background(), gen.RecentReportExistsByScannerParams{
		SourceScannerID: uuid.NullUUID{UUID: scannerID, Valid: true},
		Since:           since,
	})
	return exists, err
}

func (r *ScanReportRepo) FindByBaseAndCoordinates(baseID int, x int, y int) ([]*domain.SectorScanReport, error) {
	rows, err := r.q.ListScanReportsByBaseAndCoordinates(context.Background(), gen.ListScanReportsByBaseAndCoordinatesParams{BaseID: int64(baseID), SectorX: int32(x), SectorY: int32(y)})
	if err != nil {
		return nil, err
	}
	out := make([]*domain.SectorScanReport, 0, len(rows))
	for _, row := range rows {
		out = append(out, mappers.ScanReportFromDB(row))
	}
	return out, nil
}

func (r *ScanReportRepo) GetLatestScansByBase(baseID int) ([]*domain.SectorScanReport, error) {
	rows, err := r.q.GetLatestScanReportsByBase(context.Background(), int64(baseID))
	if err != nil {
		return nil, err
	}
	out := make([]*domain.SectorScanReport, 0, len(rows))
	for _, row := range rows {
		out = append(out, mappers.ScanReportFromDB(row))
	}
	return out, nil
}

func (r *ScanReportRepo) Delete(id int) error {
	return r.q.DeleteScanReport(context.Background(), int64(id))
}

// FindByBaseWithinArea provides a naive in-memory filtering implementation over the latest scans.
// For production efficiency, add a dedicated SQL query joining sectors with coordinate filtering.
func (r *ScanReportRepo) FindByBaseWithinArea(baseID int, centerX int, centerY int, radius int) ([]*domain.SectorScanReport, error) {
	if radius < 0 {
		return nil, errors.New("radius must be non-negative")
	}
	rows, err := r.q.ListScanReportsByBaseWithinArea(context.Background(), gen.ListScanReportsByBaseWithinAreaParams{
		BaseID:  int64(baseID),
		CenterX: int32(centerX),
		CenterY: int32(centerY),
		Radius:  int32(radius),
	})
	if err != nil {
		return nil, err
	}
	out := make([]*domain.SectorScanReport, 0, len(rows))
	for _, row := range rows {
		out = append(out, mappers.ScanReportFromDB(row))
	}
	return out, nil
}
