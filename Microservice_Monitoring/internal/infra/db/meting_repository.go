package db

import (
	"context"
	"database/sql"
	"fmt"
	"monitoring/internal/domain/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresMetingRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresMetingRepository is de constructor
func NewPostgresMetingRepository(pool *pgxpool.Pool) *PostgresMetingRepository {
	return &PostgresMetingRepository{pool: pool}
}

// Save slaat een meting op in de TimescaleDB hypertable en retourneert het volledige opgeslagen record
func (r *PostgresMetingRepository) Save(ctx context.Context, m models.Meting) (models.Meting, error) {
	query := `
		INSERT INTO meting (time, sensor_id, kunstwerk_id, waarde, is_handmatig, inspectie_id, afgehandeld)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING time, id, sensor_id, kunstwerk_id, waarde, is_handmatig, inspectie_id, afgehandeld
	`

	var (
		outTime        sql.NullTime
		outID          string
		outSensorID    sql.NullInt64
		outKunstwerkID int64
		outWaarde      float64
		outIsHandmatig bool
		outInspectieID sql.NullString
		outAfgehandeld bool
	)

	err := r.pool.QueryRow(ctx, query,
		m.Time, m.SensorID, m.KunstwerkID, m.Waarde, m.IsHandmatig, m.InspectieID, m.Afgehandeld,
	).Scan(&outTime, &outID, &outSensorID, &outKunstwerkID, &outWaarde, &outIsHandmatig, &outInspectieID, &outAfgehandeld)

	if err != nil {
		return models.Meting{}, fmt.Errorf("fout bij opslaan meting: %w", err)
	}

	saved := m
	if outTime.Valid {
		saved.Time = outTime.Time
	}
	if parsed, perr := uuid.Parse(outID); perr == nil {
		saved.ID = parsed
	}
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
	saved.Afgehandeld = outAfgehandeld

	return saved, nil
}
