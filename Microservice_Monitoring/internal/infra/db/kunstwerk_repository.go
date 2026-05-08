package db

import (
	"context"
	"errors"
	"fmt"

	"monitoring/internal/domain/models" // Pas aan naar jouw pad

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresConfiguratieRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresConfiguratieRepository(pool *pgxpool.Pool) *PostgresConfiguratieRepository {
	return &PostgresConfiguratieRepository{pool: pool}
}

// GetBySensorID haalt de specifieke bedrijfsregels voor één sensor op
func (r *PostgresConfiguratieRepository) GetBySensorID(ctx context.Context, sensorID int64) (models.SensorConfiguratie, error) {
	query := `
		SELECT id, sensor_id, min_value, max_value, marge_percentage
		FROM sensorconfiguratie
		WHERE sensor_id = $1
	`

	var config models.SensorConfiguratie

	// QueryRow voert de query uit en Scan koppelt de database-kolommen aan je Go-struct
	err := r.pool.QueryRow(ctx, query, sensorID).Scan(
		&config.ID,
		&config.SensorID,
		&config.MinValue,
		&config.MaxValue,
		&config.MargePercentage,
	)

	// Foutafhandeling: Wat als er geen configuratie is?
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Specifieke error zodat je Analyse Service weet dat hij default gedrag moet vertonen
			return models.SensorConfiguratie{}, fmt.Errorf("geen configuratie gevonden voor sensor %s", sensorID)
		}
		// Een andere database fout (bijv. verbinding weg)
		return models.SensorConfiguratie{}, fmt.Errorf("fout bij ophalen configuratie: %w", err)
	}

	return config, nil
}
