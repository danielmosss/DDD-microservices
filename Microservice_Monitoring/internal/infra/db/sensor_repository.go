package db

import (
	"context"
	"database/sql"
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

func (r *PostgresSensorRepository) GetSensorTypes(ctx context.Context) ([]models.SensorType, error) {
	query := `
		SELECT id, naam, COALESCE(eenheid, ''), drempel_is_range
		FROM sensortype
		ORDER BY naam ASC, id ASC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("fout bij ophalen sensortypes: %w", err)
	}
	defer rows.Close()

	sensorTypes := make([]models.SensorType, 0)
	for rows.Next() {
		var sensorType models.SensorType
		if err := rows.Scan(&sensorType.ID, &sensorType.Naam, &sensorType.Eenheid, &sensorType.DrempelIsRange); err != nil {
			return nil, fmt.Errorf("fout bij lezen sensortype: %w", err)
		}
		sensorTypes = append(sensorTypes, sensorType)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("fout bij itereren sensortypes: %w", err)
	}

	return sensorTypes, nil
}

func (r *PostgresSensorRepository) GetSensorTypeByID(ctx context.Context, sensorTypeID int64) (models.SensorType, error) {
	query := `
		SELECT id, naam, COALESCE(eenheid, ''), drempel_is_range
		FROM sensortype
		WHERE id = $1
	`

	var sensorType models.SensorType
	if err := r.pool.QueryRow(ctx, query, sensorTypeID).Scan(
		&sensorType.ID,
		&sensorType.Naam,
		&sensorType.Eenheid,
		&sensorType.DrempelIsRange,
	); err != nil {
		return models.SensorType{}, err
	}

	return sensorType, nil
}

func (r *PostgresSensorRepository) GetConfiguratieBronnen(ctx context.Context, kunstwerkID int64, sensorTypeID *int64) ([]models.SensorConfiguratieBron, error) {
	query := `
		SELECT
			s.id,
			s.onderdeel_id,
			o.naam,
			st.id,
			st.naam,
			st.eenheid,
			st.drempel_is_range,
			sc.id,
			sc.sensor_id,
			sc.min_value,
			sc.max_value,
			sc.marge_percentage
		FROM sensor s
		JOIN sensortype st ON s.sensortype_id = st.id
		JOIN sensorconfiguratie sc ON s.id = sc.sensor_id
		LEFT JOIN onderdelen o ON s.onderdeel_id = o.id
		WHERE s.kunstwerk_id = $1
		  AND s.deleted = false
		  AND ($2::bigint IS NULL OR st.id = $2)
		ORDER BY st.naam ASC, s.id ASC
	`

	rows, err := r.pool.Query(ctx, query, kunstwerkID, sensorTypeID)
	if err != nil {
		return nil, fmt.Errorf("fout bij ophalen configuratiebronnen: %w", err)
	}
	defer rows.Close()

	bronnen := make([]models.SensorConfiguratieBron, 0)
	for rows.Next() {
		var bron models.SensorConfiguratieBron
		var onderdeelID sql.NullInt64
		var onderdeelNaam sql.NullString
		if err := rows.Scan(
			&bron.SensorID,
			&onderdeelID,
			&onderdeelNaam,
			&bron.SensorType.ID,
			&bron.SensorType.Naam,
			&bron.SensorType.Eenheid,
			&bron.SensorType.DrempelIsRange,
			&bron.SensorConfiguratie.ID,
			&bron.SensorConfiguratie.SensorID,
			&bron.SensorConfiguratie.MinValue,
			&bron.SensorConfiguratie.MaxValue,
			&bron.SensorConfiguratie.MargePercentage,
		); err != nil {
			return nil, fmt.Errorf("fout bij lezen configuratiebron: %w", err)
		}
		if onderdeelID.Valid {
			value := onderdeelID.Int64
			bron.OnderdeelID = &value
		}
		if onderdeelNaam.Valid {
			value := onderdeelNaam.String
			bron.OnderdeelNaam = &value
		}
		bronnen = append(bronnen, bron)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("fout bij itereren configuratiebronnen: %w", err)
	}

	return bronnen, nil
}

func (r *PostgresSensorRepository) SensorBelongsToKunstwerk(ctx context.Context, sensorID int64, kunstwerkID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM sensor WHERE id = $1 AND kunstwerk_id = $2)`

	var exists bool
	if err := r.pool.QueryRow(ctx, query, sensorID, kunstwerkID).Scan(&exists); err != nil {
		return false, fmt.Errorf("fout bij controleren sensor: %w", err)
	}

	return exists, nil
}

func (r *PostgresSensorRepository) SensorExists(ctx context.Context, sensorID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM sensor WHERE id = $1)`

	var exists bool
	if err := r.pool.QueryRow(ctx, query, sensorID).Scan(&exists); err != nil {
		return false, fmt.Errorf("fout bij controleren sensor: %w", err)
	}

	return exists, nil
}

func (r *PostgresSensorRepository) DeleteSensor(ctx context.Context, sensorID int64, preserveSensorData bool) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("fout bij starten transactie: %w", err)
	}
	defer tx.Rollback(ctx)

	if preserveSensorData {
		commandTag, err := tx.Exec(ctx, `UPDATE sensor SET deleted = true WHERE id = $1`, sensorID)
		if err != nil {
			return fmt.Errorf("fout bij soft delete sensor: %w", err)
		}
		if commandTag.RowsAffected() == 0 {
			return pgx.ErrNoRows
		}
	} else {
		if _, err := tx.Exec(ctx, `DELETE FROM afwijking WHERE sensor_id = $1`, sensorID); err != nil {
			return fmt.Errorf("fout bij verwijderen afwijkingen: %w", err)
		}
		if _, err := tx.Exec(ctx, `DELETE FROM meting WHERE sensor_id = $1`, sensorID); err != nil {
			return fmt.Errorf("fout bij verwijderen metingen: %w", err)
		}
		if _, err := tx.Exec(ctx, `DELETE FROM sensorconfiguratie WHERE sensor_id = $1`, sensorID); err != nil {
			return fmt.Errorf("fout bij verwijderen sensorconfiguratie: %w", err)
		}
		commandTag, err := tx.Exec(ctx, `DELETE FROM sensor WHERE id = $1`, sensorID)
		if err != nil {
			return fmt.Errorf("fout bij verwijderen sensor: %w", err)
		}
		if commandTag.RowsAffected() == 0 {
			return pgx.ErrNoRows
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("fout bij bevestigen delete transactie: %w", err)
	}

	return nil
}

func (r *PostgresSensorRepository) CreateSensorWithConfiguratie(ctx context.Context, kunstwerkID int64, onderdeelID int64, request models.CreateSensorRequest, config models.UpdateSensorConfiguratieRequest) (models.Sensor, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return models.Sensor{}, fmt.Errorf("fout bij starten transactie: %w", err)
	}
	defer tx.Rollback(ctx)

	sensorQuery := `
		INSERT INTO sensor (kunstwerk_id, onderdeel_id, geolocation, sensortype_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, kunstwerk_id, onderdeel_id, geolocation, sensortype_id, last_analyzed_meting_id
	`

	var sensor models.Sensor
	var onderdeelIDValue sql.NullInt64
	var geolocationValue sql.NullString
	var lastAnalyzedMetingIDValue sql.NullInt64
	if err := tx.QueryRow(ctx, sensorQuery, kunstwerkID, onderdeelID, request.Geolocation, request.SensorTypeID).Scan(
		&sensor.ID,
		&sensor.KunstwerkID,
		&onderdeelIDValue,
		&geolocationValue,
		&sensor.SensorTypeID,
		&lastAnalyzedMetingIDValue,
	); err != nil {
		return models.Sensor{}, fmt.Errorf("fout bij aanmaken sensor: %w", err)
	}
	if onderdeelIDValue.Valid {
		value := onderdeelIDValue.Int64
		sensor.OnderdeelID = &value
	}
	if geolocationValue.Valid {
		value := geolocationValue.String
		sensor.Geolocation = &value
	}
	if lastAnalyzedMetingIDValue.Valid {
		value := lastAnalyzedMetingIDValue.Int64
		sensor.LastAnalyzedMetingID = &value
	}

	configQuery := `
		INSERT INTO sensorconfiguratie (sensor_id, min_value, max_value, marge_percentage)
		VALUES ($1, $2, $3, $4)
		RETURNING id, sensor_id, min_value, max_value, marge_percentage
	`

	if err := tx.QueryRow(ctx, configQuery, sensor.ID, config.MinValue, config.MaxValue, config.MargePercentage).Scan(
		&sensor.SensorConfiguratie.ID,
		&sensor.SensorConfiguratie.SensorID,
		&sensor.SensorConfiguratie.MinValue,
		&sensor.SensorConfiguratie.MaxValue,
		&sensor.SensorConfiguratie.MargePercentage,
	); err != nil {
		return models.Sensor{}, fmt.Errorf("fout bij aanmaken sensorconfiguratie: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return models.Sensor{}, fmt.Errorf("fout bij bevestigen transactie: %w", err)
	}

	return sensor, nil
}

func (r *PostgresSensorRepository) UpdateSensorConfiguratie(ctx context.Context, sensorID int64, request models.UpdateSensorConfiguratieRequest) (models.SensorConfiguratie, error) {
	query := `
		INSERT INTO sensorconfiguratie (sensor_id, min_value, max_value, marge_percentage)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (sensor_id) DO UPDATE SET
			min_value = EXCLUDED.min_value,
			max_value = EXCLUDED.max_value,
			marge_percentage = EXCLUDED.marge_percentage
		RETURNING id, sensor_id, min_value, max_value, marge_percentage
	`

	var config models.SensorConfiguratie
	if err := r.pool.QueryRow(ctx, query, sensorID, request.MinValue, request.MaxValue, request.MargePercentage).Scan(
		&config.ID,
		&config.SensorID,
		&config.MinValue,
		&config.MaxValue,
		&config.MargePercentage,
	); err != nil {
		return models.SensorConfiguratie{}, fmt.Errorf("fout bij opslaan sensorconfiguratie: %w", err)
	}

	return config, nil
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
	'deleted', s.deleted,
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
