package db

import (
	"context"
	"database/sql"
	"fmt"
	"monitoring/internal/domain/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresMetingRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresMetingRepository is de constructor
func NewPostgresMetingRepository(pool *pgxpool.Pool) *PostgresMetingRepository {
	return &PostgresMetingRepository{pool: pool}
}

func (r *PostgresMetingRepository) GetMetingenByKunstwerkID(ctx context.Context, kunstwerkID int64, limit int, offset int) ([]models.Meting, int64, error) {
	query := `
		SELECT time, id, sensor_id, kunstwerk_id, waarde, is_handmatig, inspectie_id
		FROM meting
		WHERE kunstwerk_id = $1
		ORDER BY time DESC, id DESC
		LIMIT $2
		OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, kunstwerkID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("fout bij ophalen metingen: %w", err)
	}
	defer rows.Close()

	metingen := make([]models.Meting, 0)
	for rows.Next() {
		var (
			item           models.Meting
			outTime        sql.NullTime
			outSensorID    sql.NullInt64
			outInspectieID sql.NullString
		)

		if err := rows.Scan(
			&outTime,
			&item.ID,
			&outSensorID,
			&item.KunstwerkID,
			&item.Waarde,
			&item.IsHandmatig,
			&outInspectieID,
		); err != nil {
			return nil, 0, fmt.Errorf("fout bij lezen meting: %w", err)
		}

		if outTime.Valid {
			item.Time = outTime.Time
		}
		if outSensorID.Valid {
			value := outSensorID.Int64
			item.SensorID = &value
		} else {
			item.SensorID = nil
		}
		if outInspectieID.Valid {
			value := outInspectieID.String
			item.InspectieID = &value
		} else {
			item.InspectieID = nil
		}

		metingen = append(metingen, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("fout bij itereren metingen: %w", err)
	}

	total, err := r.countMetingenByKunstwerkID(ctx, kunstwerkID)
	if err != nil {
		return nil, 0, fmt.Errorf("fout bij tellen metingen: %w", err)
	}

	return metingen, total, nil

func (r *PostgresMetingRepository) countMetingenByKunstwerkID(ctx context.Context, kunstwerkID int64) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM meting
		WHERE kunstwerk_id = $1
	`

	var total int64
	if err := r.pool.QueryRow(ctx, query, kunstwerkID).Scan(&total); err != nil {
		return 0, fmt.Errorf("fout bij tellen metingen: %w", err)
	}

	return total, nil
}

func (r *PostgresMetingRepository) GetRecentMetingPerSensorByKunstwerkID(ctx context.Context, kunstwerkID int64) ([]models.Meting, error) {
	query := `
		SELECT DISTINCT ON (sensor_id) time, id, sensor_id, kunstwerk_id, waarde, is_handmatig, inspectie_id
		FROM meting
		WHERE kunstwerk_id = $1 AND sensor_id IS NOT NULL
		ORDER BY sensor_id, time DESC, id DESC
	`

	rows, err := r.pool.Query(ctx, query, kunstwerkID)
	if err != nil {
		return nil, fmt.Errorf("fout bij ophalen recente metingen: %w", err)
	}
	defer rows.Close()

	metingen := make([]models.Meting, 0)
	for rows.Next() {
		var (
			item           models.Meting
			outTime        sql.NullTime
			outSensorID    sql.NullInt64
			outInspectieID sql.NullString
		)

		if err := rows.Scan(
			&outTime,
			&item.ID,
			&outSensorID,
			&item.KunstwerkID,
			&item.Waarde,
			&item.IsHandmatig,
			&outInspectieID,
		); err != nil {
			return nil, fmt.Errorf("fout bij lezen recente meting: %w", err)
		}

		if outTime.Valid {
			item.Time = outTime.Time
		}
		if outSensorID.Valid {
			value := outSensorID.Int64
			item.SensorID = &value
		} else {
			item.SensorID = nil
		}
		if outInspectieID.Valid {
			value := outInspectieID.String
			item.InspectieID = &value
		} else {
			item.InspectieID = nil
		}

		metingen = append(metingen, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("fout bij itereren recente metingen: %w", err)
	}

	return metingen, nil
}

// Save slaat een meting op in de TimescaleDB hypertable en retourneert het volledige opgeslagen record
func (r *PostgresMetingRepository) Save(ctx context.Context, m models.Meting, returnObject bool) (models.Meting, error) {
	query := `
		INSERT INTO meting (time, sensor_id, kunstwerk_id, waarde, is_handmatig, inspectie_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING time, id, sensor_id, kunstwerk_id, waarde, is_handmatig, inspectie_id
	`

	var (
		outTime        sql.NullTime
		outID          int64
		outSensorID    sql.NullInt64
		outKunstwerkID int64
		outWaarde      float64
		outIsHandmatig bool
		outInspectieID sql.NullString
	)

	err := r.pool.QueryRow(ctx, query,
		m.Time, m.SensorID, m.KunstwerkID, m.Waarde, m.IsHandmatig, m.InspectieID,
	).Scan(&outTime, &outID, &outSensorID, &outKunstwerkID, &outWaarde, &outIsHandmatig, &outInspectieID)

	if err != nil {
		return models.Meting{}, fmt.Errorf("fout bij opslaan meting: %w", err)
	}

	if !returnObject {
		return models.Meting{}, nil
	}

	saved := m
	if outTime.Valid {
		saved.Time = outTime.Time
	}
	saved.ID = outID
	if outSensorID.Valid {
		v := outSensorID.Int64
		saved.SensorID = &v
	} else {
		saved.SensorID = nil
	}
	saved.KunstwerkID = outKunstwerkID
	saved.Waarde = outWaarde
	saved.IsHandmatig = outIsHandmatig
	if outInspectieID.Valid {
		v := outInspectieID.String
		saved.InspectieID = &v
	} else {
		saved.InspectieID = nil
	}

	return saved, nil
}
