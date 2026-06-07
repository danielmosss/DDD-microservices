package db

import (
	"context"
	"database/sql"
	"fmt"
	"monitoring/internal/domain/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresAfwijkingRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresAfwijkingRepository(pool *pgxpool.Pool) *PostgresAfwijkingRepository {
	return &PostgresAfwijkingRepository{pool: pool}
}

func (r *PostgresAfwijkingRepository) GetAfwijkingByKunstwerkID(ctx context.Context, kunstwerkID int64, limit int, offset int) ([]models.Afwijking, int64, error) {
	query := `
		SELECT id, meting_id, meting_time, kunstwerk_id, sensor_id, time, norm_min_waarde, norm_max_waarde, norm_marge_percentage, gemeten_waarde, is_warning
		FROM afwijking
		WHERE kunstwerk_id = $1
		ORDER BY time DESC, id DESC
		LIMIT $2
		OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, kunstwerkID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("fout bij ophalen afwijkingen: %w", err)
	}
	defer rows.Close()

	afwijkingen := make([]models.Afwijking, 0)
	for rows.Next() {
		var (
			item                   models.Afwijking
			outMetingTime          sql.NullTime
			outSensorID            sql.NullInt64
			outTime                sql.NullTime
			outNormMaxWaarde       sql.NullFloat64
			outNormMargePercentage sql.NullFloat64
		)

		if err := rows.Scan(
			&item.ID,
			&item.MetingID,
			&outMetingTime,
			&item.KunstwerkID,
			&outSensorID,
			&outTime,
			&item.NormMinWaarde,
			&outNormMaxWaarde,
			&outNormMargePercentage,
			&item.GemetenWaarde,
			&item.IsWarning,
		); err != nil {
			return nil, 0, fmt.Errorf("fout bij lezen afwijking: %w", err)
		}

		if outMetingTime.Valid {
			item.MetingTime = outMetingTime.Time
		}
		if outSensorID.Valid {
			value := outSensorID.Int64
			item.SensorID = &value
		} else {
			item.SensorID = nil
		}
		if outTime.Valid {
			item.Time = outTime.Time
		}
		if outNormMaxWaarde.Valid {
			item.NormMaxWaarde = outNormMaxWaarde.Float64
		}
		if outNormMargePercentage.Valid {
			value := outNormMargePercentage.Float64
			item.NormMargePercentage = &value
		} else {
			item.NormMargePercentage = nil
		}

		afwijkingen = append(afwijkingen, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("fout bij itereren afwijkingen: %w", err)
	}

	total, err := r.countByKunstwerkID(ctx, kunstwerkID)
	if err != nil {
		return nil, 0, err
	}

	return afwijkingen, total, nil
}

func (r *PostgresAfwijkingRepository) countAfwijkingenByKunstwerkID(ctx context.Context, kunstwerkID int64) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM afwijking
		WHERE kunstwerk_id = $1
	`

	var total int64
	if err := r.pool.QueryRow(ctx, query, kunstwerkID).Scan(&total); err != nil {
		return 0, fmt.Errorf("fout bij tellen afwijkingen: %w", err)
	}

	return total, nil
}

func (r *PostgresAfwijkingRepository) Save(ctx context.Context, m models.Afwijking) (models.Afwijking, error) {
	query := `
		INSERT INTO afwijking (meting_id, meting_time, kunstwerk_id, sensor_id, time, norm_min_waarde, norm_max_waarde, gemeten_waarde, is_warning)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, meting_id, meting_time, kunstwerk_id, sensor_id, time, norm_min_waarde, norm_max_waarde, gemeten_waarde, is_warning
	`

	var (
		outId            int64
		outMetingId      int64
		outMetingTime    sql.NullTime
		outKunstwerkID   int64
		outSensorID      sql.NullInt64
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
	saved.MetingID = outMetingId
	if outMetingTime.Valid {
		saved.MetingTime = outMetingTime.Time
	}
	saved.KunstwerkID = outKunstwerkID
	if outSensorID.Valid {
		v := outSensorID.Int64
		saved.SensorID = &v
	} else {
		saved.SensorID = nil
	}
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
