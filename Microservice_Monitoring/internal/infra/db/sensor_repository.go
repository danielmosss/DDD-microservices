package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"monitoring/internal/domain/models" // Pas aan naar jouw pad

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresSensorRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresSensorRepository(pool *pgxpool.Pool) *PostgresSensorRepository {
	return &PostgresSensorRepository{pool: pool}
}

func (r *PostgresSensorRepository) GetBySensorID(ctx context.Context, sensorID int64) (models.SensorConfiguratie, error) {
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
			return models.SensorConfiguratie{}, fmt.Errorf("geen configuratie gevonden voor sensor %d", sensorID)
		}
		// Een andere database fout (bijv. verbinding weg)
		return models.SensorConfiguratie{}, fmt.Errorf("fout bij ophalen configuratie: %w", err)
	}

	return config, nil
}

func (r *PostgresSensorRepository) GetSensorAndLastMetingAfwijkingFromSensorId(ctx context.Context, kunstwerkId int64, sensorId int64) (models.SensorDetailResponse, error) {
	query := `
SELECT json_build_object(
    'id', s.id,
    'sensorType', row_to_json(st),
    'sensorConfiguratie', row_to_json(sc),
    'laatsteMeting', CASE WHEN m.id IS NOT NULL THEN row_to_json(m) ELSE NULL END,
    'afwijking', CASE WHEN a.sensor_id IS NOT NULL THEN row_to_json(a) ELSE NULL END
)
FROM sensor s
JOIN sensorconfiguratie sc on s.id = sc.sensor_id
JOIN sensortype st on s.sensortype_id = st.id
LEFT JOIN meting m on s.id = m.sensor_id and s.last_analyzed_meting_id = m.id
LEFT JOIN afwijking a on s.id = a.sensor_id and s.last_analyzed_meting_id = a.meting_id
WHERE s.id = @sensor_id and s.kunstwerk_id = @kunstwerk_id
`

	var SensorDetailResponse models.SensorDetailResponse
	var jsondata []byte

	args := pgx.NamedArgs{
		"sensor_id":    sensorId,
		"kunstwerk_id": kunstwerkId,
	}

	err := r.pool.QueryRow(ctx, query, args).Scan(&jsondata)
	if err != nil {
		return models.SensorDetailResponse{}, err
	}

	err = json.Unmarshal(jsondata, &SensorDetailResponse)
	if err != nil {
		return models.SensorDetailResponse{}, err
	}

	SensorDetailResponse.Status = models.StatusHealthy
	if SensorDetailResponse.Afwijking != nil {
		if SensorDetailResponse.Afwijking.IsWarning {
			SensorDetailResponse.Status = models.StatusWarning
		} else {
			SensorDetailResponse.Status = models.StatusCritical
		}
	}

	return SensorDetailResponse, nil
}
