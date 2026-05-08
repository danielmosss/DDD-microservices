package db

import (
	"context"
	"database/sql"
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

func (r *PostgresKunstwerkRepository) GetKunstwerkMetType(ctx context.Context, kunstwerkID int64) (models.KunstwerkDetail, error) {
	query := `
        SELECT 
            k.id,
			k.beheeridentifier,
			k.naam,
			k.geolocation,
			COALESCE(k.kunstwerktype_id, 0),
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
	var kunstwerkBeschrijving sql.NullString
	var kunstwerkTypeBeschrijving sql.NullString

	// Voer de query uit en map direct naar de velden in je geneste structs
	err := r.pool.QueryRow(ctx, query, kunstwerkID).Scan(
		&result.Kunstwerk.ID,
		&result.Kunstwerk.BeheerIdentifier,
		&result.Kunstwerk.Naam,
		&kunstwerkGeolocation,
		&result.Kunstwerk.KunstwerkTypeID,
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