package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"monitoring/internal/domain/models" // Pas aan naar jouw pad

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresKunstwerkRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresKunstwerkRepository(pool *pgxpool.Pool) *PostgresKunstwerkRepository {
	return &PostgresKunstwerkRepository{pool: pool}
}

func (r *PostgresKunstwerkRepository) KunstwerkExists(ctx context.Context, kunstwerkID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM kunstwerk WHERE id = $1 AND deleted = false)`

	var exists bool
	if err := r.pool.QueryRow(ctx, query, kunstwerkID).Scan(&exists); err != nil {
		return false, fmt.Errorf("fout bij controleren kunstwerk: %w", err)
	}

	return exists, nil
}

func (r *PostgresKunstwerkRepository) OnderdeelBelongsToKunstwerk(ctx context.Context, onderdeelID int64, kunstwerkID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM onderdelen WHERE id = $1 AND kunstwerk_id = $2 AND deleted = false)`

	var exists bool
	if err := r.pool.QueryRow(ctx, query, onderdeelID, kunstwerkID).Scan(&exists); err != nil {
		return false, fmt.Errorf("fout bij controleren onderdeel: %w", err)
	}

	return exists, nil
}

func (r *PostgresKunstwerkRepository) CreateOnderdeel(ctx context.Context, kunstwerkID int64, request models.CreateOnderdeelRequest) (models.KunstwerkOnderdeel, error) {
	if request.ParentOnderdeelID != nil {
		exists, err := r.OnderdeelBelongsToKunstwerk(ctx, *request.ParentOnderdeelID, kunstwerkID)
		if err != nil {
			return models.KunstwerkOnderdeel{}, err
		}
		if !exists {
			return models.KunstwerkOnderdeel{}, pgx.ErrNoRows
		}
	}

	query := `
		INSERT INTO onderdelen (kunstwerk_id, naam, parent_id)
		VALUES ($1, $2, $3)
		RETURNING id, kunstwerk_id, naam, parent_id
	`

	var onderdeel models.KunstwerkOnderdeel
	if err := r.pool.QueryRow(ctx, query, kunstwerkID, request.Naam, request.ParentOnderdeelID).Scan(
		&onderdeel.ID,
		&onderdeel.KunstwerkId,
		&onderdeel.Naam,
		&onderdeel.ParentId,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.KunstwerkOnderdeel{}, err
		}
		return models.KunstwerkOnderdeel{}, fmt.Errorf("fout bij aanmaken onderdeel: %w", err)
	}

	return onderdeel, nil
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
		SELECT sensor.id, kunstwerk_id, onderdeel_id, geolocation, sensortype_id, last_analyzed_meting_id, sensor.deleted, sc.*
		FROM sensor
		LEFT JOIN sensorconfiguratie sc ON sensor.id = sc.sensor_id
		WHERE kunstwerk_id = $1
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
			&item.Deleted,
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

func (r *PostgresKunstwerkRepository) GetKunstwerkOnderdelen(ctx context.Context, kunstwerkId int64) ([]models.KunstwerkOnderdeel, error) {
	query := `
SELECT o.id, o.naam, o.parent_id, o.deleted
FROM onderdelen o
WHERE kunstwerk_id = $1;
	`

	var onderdelen []models.KunstwerkOnderdeel
	rows, err := r.pool.Query(ctx, query, kunstwerkId)
	if err != nil {
		return nil, fmt.Errorf("fout bij ophalen onderdelen: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var onderdeel models.KunstwerkOnderdeel
		err := rows.Scan(&onderdeel.ID, &onderdeel.Naam, &onderdeel.ParentId, &onderdeel.Deleted)
		if err != nil {
			return nil, err
		}
		onderdeel.KunstwerkId = kunstwerkId
		onderdelen = append(onderdelen, onderdeel)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("fout bij itereren sensoren: %w", err)
	}
	return onderdelen, nil
}

func (r *PostgresKunstwerkRepository) GetKunstwerkOnderdelenWithSensorIDs(ctx context.Context, kunstwerkId int64) ([]models.KunstwerkOnderdeelMetSensor, error) {
	query := `
SELECT
    o.id,
    o.naam,
    o.parent_id,
	o.deleted,
    COALESCE(array_remove(array_agg(s.id), NULL), '{}') AS sensor_ids
FROM onderdelen o
LEFT JOIN sensor s ON o.id = s.onderdeel_id
WHERE o.kunstwerk_id = $1
GROUP BY o.id, o.naam, o.parent_id, o.deleted;
	`

	var onderdelen []models.KunstwerkOnderdeelMetSensor
	rows, err := r.pool.Query(ctx, query, kunstwerkId)
	if err != nil {
		return nil, fmt.Errorf("fout bij ophalen onderdelen met sensoren: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var onderdeel models.KunstwerkOnderdeelMetSensor
		err := rows.Scan(&onderdeel.ID, &onderdeel.Naam, &onderdeel.ParentId, &onderdeel.Deleted, &onderdeel.SensorIds)
		if err != nil {
			return nil, err
		}
		onderdeel.KunstwerkId = kunstwerkId
		onderdelen = append(onderdelen, onderdeel)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("fout bij itereren sensoren: %w", err)
	}
	return onderdelen, nil
}

func (r *PostgresKunstwerkRepository) DeleteOnderdeelTree(ctx context.Context, kunstwerkID int64, onderdeelID int64, preserveSensorData bool) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("fout bij starten transactie: %w", err)
	}
	defer tx.Rollback(ctx)

	subtreeQuery := `
		WITH RECURSIVE subtree AS (
			SELECT id
			FROM onderdelen
			WHERE id = $1 AND kunstwerk_id = $2
			UNION ALL
			SELECT child.id
			FROM onderdelen child
			JOIN subtree parent ON child.parent_id = parent.id
			WHERE child.kunstwerk_id = $2
		)
		SELECT id FROM subtree
	`

	rows, err := tx.Query(ctx, subtreeQuery, onderdeelID, kunstwerkID)
	if err != nil {
		return fmt.Errorf("fout bij ophalen onderdeelboom: %w", err)
	}

	onderdeelIDs := make([]int64, 0)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return fmt.Errorf("fout bij lezen onderdeelboom: %w", err)
		}
		onderdeelIDs = append(onderdeelIDs, id)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("fout bij itereren onderdeelboom: %w", err)
	}
	rows.Close()

	if len(onderdeelIDs) == 0 {
		return pgx.ErrNoRows
	}

	if preserveSensorData {
		if _, err := tx.Exec(ctx, `UPDATE sensor SET deleted = true WHERE kunstwerk_id = $1 AND onderdeel_id = ANY($2)`, kunstwerkID, onderdeelIDs); err != nil {
			return fmt.Errorf("fout bij soft delete sensoren: %w", err)
		}
		if _, err := tx.Exec(ctx, `UPDATE onderdelen SET deleted = true WHERE kunstwerk_id = $1 AND id = ANY($2)`, kunstwerkID, onderdeelIDs); err != nil {
			return fmt.Errorf("fout bij soft delete onderdelen: %w", err)
		}
	} else {
		if _, err := tx.Exec(ctx, `DELETE FROM afwijking WHERE sensor_id IN (SELECT id FROM sensor WHERE kunstwerk_id = $1 AND onderdeel_id = ANY($2))`, kunstwerkID, onderdeelIDs); err != nil {
			return fmt.Errorf("fout bij verwijderen afwijkingen: %w", err)
		}
		if _, err := tx.Exec(ctx, `DELETE FROM meting WHERE sensor_id IN (SELECT id FROM sensor WHERE kunstwerk_id = $1 AND onderdeel_id = ANY($2))`, kunstwerkID, onderdeelIDs); err != nil {
			return fmt.Errorf("fout bij verwijderen metingen: %w", err)
		}
		if _, err := tx.Exec(ctx, `DELETE FROM sensorconfiguratie WHERE sensor_id IN (SELECT id FROM sensor WHERE kunstwerk_id = $1 AND onderdeel_id = ANY($2))`, kunstwerkID, onderdeelIDs); err != nil {
			return fmt.Errorf("fout bij verwijderen sensorconfiguratie: %w", err)
		}
		if _, err := tx.Exec(ctx, `DELETE FROM sensor WHERE kunstwerk_id = $1 AND onderdeel_id = ANY($2)`, kunstwerkID, onderdeelIDs); err != nil {
			return fmt.Errorf("fout bij verwijderen sensoren: %w", err)
		}
		if _, err := tx.Exec(ctx, `DELETE FROM onderdelen WHERE kunstwerk_id = $1 AND id = ANY($2)`, kunstwerkID, onderdeelIDs); err != nil {
			return fmt.Errorf("fout bij verwijderen onderdelen: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("fout bij bevestigen delete transactie: %w", err)
	}

	return nil
}

func (r *PostgresKunstwerkRepository) GetAantalSensoren(ctx context.Context, kunstwerkId int64) (int, error) {
	query := `
SELECT COUNT(DISTINCT s.id)
FROM sensor s
WHERE s.kunstwerk_id = $1
  AND s.deleted = false
	`

	var count int
	err := r.pool.QueryRow(ctx, query, kunstwerkId).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *PostgresKunstwerkRepository) GetAantalActieveSensoren(ctx context.Context, kunstwerkId int64) (int, error) {
	query := `
SELECT COUNT(DISTINCT s.id) AS active_sensors_with_recent_data
FROM sensor s
WHERE s.kunstwerk_id = $1
  AND s.deleted = false
  AND EXISTS (
      SELECT 1
      FROM meting m
      WHERE m.sensor_id = s.id
        AND m.time > NOW() - INTERVAL '24 hours'
  );
	`

	var count int
	err := r.pool.QueryRow(ctx, query, kunstwerkId).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *PostgresKunstwerkRepository) GetAantalAfwijkingen(ctx context.Context, kunstwerkId int64) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM afwijking a
		WHERE a.kunstwerk_id = $1 AND a.time >= NOW() - INTERVAL '24 hours'
	`

	var count int
	err := r.pool.QueryRow(ctx, query, kunstwerkId).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *PostgresKunstwerkRepository) GetAantalSensorenMetNAfwijkingen(ctx context.Context, kunstwerkId int64) (int, error) {
	query := `
		SELECT COUNT(DISTINCT sensor_id)
		FROM afwijking a
		WHERE a.kunstwerk_id = $1 AND a.time >= NOW() - INTERVAL '24 hours'
	`

	var count int
	err := r.pool.QueryRow(ctx, query, kunstwerkId).Scan(&count)
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

func (r *PostgresKunstwerkRepository) InsertUpdateKunstwerkDHU(ctx context.Context, kunstwerkId int64, DHS *models.DailyHealthSummary) error {
	query := `
INSERT INTO dailyhealthupdatecache (kunstwerkid, status, aantalsensoren, aantalactievesensoren, aantalafwijkendesensoren, aantalafwijkingen)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (kunstwerkid)
DO UPDATE SET
    status = EXCLUDED.status,
    aantalsensoren = EXCLUDED.aantalsensoren,
    aantalactievesensoren = EXCLUDED.aantalactievesensoren,
    aantalafwijkendesensoren = EXCLUDED.aantalafwijkendesensoren,
    aantalafwijkingen = EXCLUDED.aantalafwijkingen;
	`
	_, err := r.pool.Exec(ctx, query, kunstwerkId, DHS.Status, DHS.AantalSensoren, DHS.AantalActieveSensoren, DHS.AantalAfwijkendeSensoren, DHS.AantalAfwijkingen)
	return err
}

func (r *PostgresKunstwerkRepository) GetKunstwerkDHU(ctx context.Context, kunstwerkId int64) (models.DailyHealthUpdate, error) {
	query := ` 
SELECT *
FROM dailyhealthupdatecache
WHERE kunstwerkid = $1
`

	var result models.DailyHealthUpdate
	err := r.pool.QueryRow(ctx, query, kunstwerkId).Scan(
		&result.KunstwerkID,
		&result.Status,
		&result.AantalSensoren,
		&result.AantalActieveSensoren,
		&result.AantalAfwijkendeSensoren,
		&result.AantalAfwijkingen,
	)
	if err != nil {
		return models.DailyHealthUpdate{}, err
	}
	return result, nil
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
