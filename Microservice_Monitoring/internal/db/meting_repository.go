package db

import (
	"context"
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

// Save slaat een meting op in de TimescaleDB hypertable
func (r *PostgresMetingRepository) Save(ctx context.Context, m models.Meting) error {
	query := `
		INSERT INTO meting (time, sensor_id, kunstwerk_id, waarde, is_afwijking, is_handmatig, inspectie_id, afgehandeld)
		VALUES ($1, $2, $3, $4, $5, $6, $7, false)
	`

	err := r.pool.QueryRow(ctx, query,
		m.Time, m.SensorID, m.KunstwerkID, m.Waarde,
		m.IsAfwijking, m.IsHandmatig, m.InspectieID, m.Afgehandeld,
	)

	if err != nil {
		return fmt.Errorf("fout bij opslaan meting: %w", err)
	}

	return nil
}
