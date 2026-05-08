package db

import (
	"context"
	"database/sql"
	"fmt"
	"monitoring/internal/domain/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresAfwijkingRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresAfwijkingRepository(pool *pgxpool.Pool) *PostgresAfwijkingRepository {
	return &PostgresAfwijkingRepository{pool: pool}
}

func (r *PostgresAfwijkingRepository) Save(ctx context.Context, m models.Afwijking) (models.Afwijking, error) {
	query := `
		INSERT INTO afwijking (meting_id, meting_time, kunstwerk_id, sensor_id, time, norm_min_waarde, norm_max_waarde, gemeten_waarde, is_warning)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, meting_id, meting_time, kunstwerk_id, sensor_id, time, norm_min_waarde, norm_max_waarde, gemeten_waarde, is_warning
	`

	var (
		outId            int64
		outMetingId      string
		outMetingTime    sql.NullTime
		outKunstwerkID   int64
		outSensorID      int64
		outTime          sql.NullTime
		outNormMinWaarde float64
		outNormMaxWaarde sql.NullFloat64
		outGemetenWaarde float64
		outIsWarning     bool
	)

	err := r.pool.QueryRow(ctx, query,
		m.MetingID, m.MetingTime, m.KunstwerkID, m.SensorID, m.Time, m.NormMinWaarde, m.NormMaxWaarde, m.GemetenWaarde, m.IsWarning,
	).Scan(&outId, &outMetingId, &outMetingTime, &outKunstwerkID, &outSensorID, &outTime, &outNormMinWaarde, &outNormMaxWaarde, &outGemetenWaarde, &outIsWarning)
	if err != nil {
		return models.Afwijking{}, fmt.Errorf("fout bij opslaan afwijking: %w", err)
	}

	saved := m
	saved.ID = outId
	if parsed, perr := uuid.Parse(outMetingId); perr == nil {
		saved.MetingID = parsed
	}
	if outMetingTime.Valid {
		saved.MetingTime = outMetingTime.Time
	}
	saved.KunstwerkID = outKunstwerkID
	saved.SensorID = outSensorID
	if outTime.Valid {
		saved.Time = outTime.Time
	}
	saved.NormMinWaarde = outNormMinWaarde
	if outNormMaxWaarde.Valid {
		saved.NormMaxWaarde = outNormMaxWaarde.Float64
	}
	saved.GemetenWaarde = outGemetenWaarde
	saved.IsWarning = outIsWarning

	return saved, nil
}
