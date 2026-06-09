package db

import (
	"context"
	"database/sql"
	"fmt"
	"monitoring/internal/domain/models" // Pas aan naar jouw pad
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresKunstwerkRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresKunstwerkRepository(pool *pgxpool.Pool) *PostgresKunstwerkRepository {
	return &PostgresKunstwerkRepository{pool: pool}
}

func (r *PostgresKunstwerkRepository) GetActieveKunstwerken(ctx context.Context) ([]models.Kunstwerk, error) {
	query := `
		SELECT id, beheeridentifier, naam, geolocation, kunstwerktype_id, beschrijving, deleted, last_send_dh_update
		FROM kunstwerk
		WHERE deleted = false
		ORDER BY naam ASC, id ASC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("fout bij ophalen kunstwerken: %w", err)
	}
	defer rows.Close()

	kunstwerken := make([]models.Kunstwerk, 0)
	for rows.Next() {
		var (
			item                  models.Kunstwerk
			kunstwerkGeolocation  sql.NullString
			kunstwerkTypeID       sql.NullInt64
			kunstwerkBeschrijving sql.NullString
			lastSendDhUpdate      sql.NullTime
		)

		if err := rows.Scan(
			&item.ID,
			&item.BeheerIdentifier,
			&item.Naam,
			&kunstwerkGeolocation,
			&kunstwerkTypeID,
			&kunstwerkBeschrijving,
			&item.Deleted,
			&lastSendDhUpdate,
		); err != nil {
			return nil, fmt.Errorf("fout bij lezen kunstwerk: %w", err)
		}

		if kunstwerkGeolocation.Valid {
			value := kunstwerkGeolocation.String
			item.Geolocation = &value
		}
		if kunstwerkTypeID.Valid {
			value := kunstwerkTypeID.Int64
			item.KunstwerkTypeID = &value
		}
		if kunstwerkBeschrijving.Valid {
			value := kunstwerkBeschrijving.String
			item.Beschrijving = &value
		}
		if lastSendDhUpdate.Valid {
			value := lastSendDhUpdate.Time
			item.LastSendDhUpdate = &value
		}

		kunstwerken = append(kunstwerken, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("fout bij itereren kunstwerken: %w", err)
	}

	return kunstwerken, nil
}

func (r *PostgresKunstwerkRepository) GetSensorenByKunstwerkID(ctx context.Context, kunstwerkID int64) ([]models.Sensor, error) {
	query := `
		SELECT sensor.id, kunstwerk_id, onderdeel_id, geolocation, sensortype_id, last_analyzed_meting_id, sc.*
		FROM sensor
		LEFT JOIN sensorconfiguratie sc ON sensor.id = sc.sensor_id
		WHERE kunstwerk_id = $1 AND deleted = false
		ORDER BY sensor.id ASC
	`

	rows, err := r.pool.Query(ctx, query, kunstwerkID)
	if err != nil {
		return nil, fmt.Errorf("fout bij ophalen sensoren: %w", err)
	}
	defer rows.Close()

	sensoren := make([]models.Sensor, 0)
	for rows.Next() {
		var (
			item                 models.Sensor
			onderdeelID          sql.NullInt64
			geolocation          sql.NullString
			lastAnalyzedMetingID sql.NullInt64
		)

		if err := rows.Scan(
			&item.ID,
			&item.KunstwerkID,
			&onderdeelID,
			&geolocation,
			&item.SensorTypeID,
			&lastAnalyzedMetingID,
			&item.SensorConfiguratie.ID,
			&item.SensorConfiguratie.SensorID,
			&item.SensorConfiguratie.MinValue,
			&item.SensorConfiguratie.MaxValue,
			&item.SensorConfiguratie.MargePercentage,
		); err != nil {
			return nil, fmt.Errorf("fout bij lezen sensor: %w", err)
		}

		if onderdeelID.Valid {
			value := onderdeelID.Int64
			item.OnderdeelID = &value
		}
		if geolocation.Valid {
			value := geolocation.String
			item.Geolocation = &value
		}
		if lastAnalyzedMetingID.Valid {
			value := lastAnalyzedMetingID.Int64
			item.LastAnalyzedMetingID = &value
		}

		sensoren = append(sensoren, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("fout bij itereren sensoren: %w", err)
	}

	return sensoren, nil
}

func (r *PostgresKunstwerkRepository) GetKunstwerkMetType(ctx context.Context, kunstwerkID int64) (models.KunstwerkDetail, error) {
	query := `
        SELECT 
            k.id,
			k.beheeridentifier,
			k.naam,
			k.geolocation,
			k.kunstwerktype_id,
			k.beschrijving,
			k.deleted,
			COALESCE(kt.id, 0),
			COALESCE(kt.naam, ''),
			kt.beschrijving
        FROM kunstwerk k
		LEFT JOIN kunstwerktype kt ON k.kunstwerktype_id = kt.id
        WHERE k.id = $1 AND k.deleted = false
    `

	var result models.KunstwerkDetail
	var kunstwerkGeolocation sql.NullString
	var kunstwerkTypeID sql.NullInt64
	var kunstwerkBeschrijving sql.NullString
	var kunstwerkTypeBeschrijving sql.NullString

	// Voer de query uit en map direct naar de velden in je geneste structs
	err := r.pool.QueryRow(ctx, query, kunstwerkID).Scan(
		&result.Kunstwerk.ID,
		&result.Kunstwerk.BeheerIdentifier,
		&result.Kunstwerk.Naam,
		&kunstwerkGeolocation,
		&kunstwerkTypeID,
		&kunstwerkBeschrijving,
		&result.Kunstwerk.Deleted,
		// En hier map je de kolommen van de gejoinde tabel:
		&result.KunstwerkType.ID,
		&result.KunstwerkType.Naam,
		&kunstwerkTypeBeschrijving,
	)

	if err != nil {
		return models.KunstwerkDetail{}, err
	}

	if kunstwerkGeolocation.Valid {
		result.Kunstwerk.Geolocation = &kunstwerkGeolocation.String
	}
	if kunstwerkTypeID.Valid {
		result.Kunstwerk.KunstwerkTypeID = &kunstwerkTypeID.Int64
	}
	if kunstwerkBeschrijving.Valid {
		result.Kunstwerk.Beschrijving = &kunstwerkBeschrijving.String
	}
	if kunstwerkTypeBeschrijving.Valid {
		result.KunstwerkType.Beschrijving = &kunstwerkTypeBeschrijving.String
	}

	return result, nil
}

func (r *PostgresKunstwerkRepository) GetAantalActieveSensoren(ctx context.Context, kunstwerkId int64) (int, error) {
	query := `
		SELECT COUNT(DISTINCT s.id)
		FROM sensor s
		WHERE s.kunstwerk_id = $1 AND s.deleted = false
	`

	var count int
	err := r.pool.QueryRow(ctx, query, kunstwerkId).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *PostgresKunstwerkRepository) GetAantalAfwijkingen(ctx context.Context, kunstwerkId int64, sinds time.Time) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM afwijking a
		WHERE a.kunstwerk_id = $1 AND a.time >= $2
	`

	var count int
	err := r.pool.QueryRow(ctx, query, kunstwerkId, sinds).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *PostgresKunstwerkRepository) GetAantalSensorenMetNAfwijkingen(ctx context.Context, kunstwerkId int64, sinds time.Time) (int, error) {
	query := `
		SELECT COUNT(DISTINCT sensor_id)
		FROM afwijking a
		WHERE a.kunstwerk_id = $1 AND a.time >= $2
	`

	var count int
	err := r.pool.QueryRow(ctx, query, kunstwerkId, sinds).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *PostgresKunstwerkRepository) UpdateKunstwerkDHupdateTime(ctx context.Context, kunstwerkId int64) error {
	query := `
		UPDATE kunstwerk
		SET last_send_dh_update = NOW()
		WHERE id = $1
	`
	_, err := r.pool.Exec(ctx, query, kunstwerkId)
	return err
}

func (r *PostgresKunstwerkRepository) GetKunstwerkenNeedingReport(ctx context.Context) ([]int64, error) {
	query := `
		SELECT id
		FROM kunstwerk
		WHERE last_send_dh_update IS NULL OR last_send_dh_update < NOW() - INTERVAL '24 hours'
		LIMIT 5
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var kunstwerkIds []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		kunstwerkIds = append(kunstwerkIds, id)
	}

	return kunstwerkIds, nil
}
