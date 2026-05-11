package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AnalysisResult struct {
	SensorID           int64
	KunstwerkID        int64
	MetingenProcessed  int64
	AfwijkingenDetected int64
	LastMetingID       int64
	Status             string
}

type AnalysisProcedureRepository struct {
	pool *pgxpool.Pool
}

func NewAnalysisProcedureRepository(pool *pgxpool.Pool) *AnalysisProcedureRepository {
	return &AnalysisProcedureRepository{pool: pool}
}

// ExecuteAnalysis runs the analyze_sensor_metingen stored procedure and returns results per sensor
func (r *AnalysisProcedureRepository) ExecuteAnalysis(ctx context.Context) ([]AnalysisResult, error) {
	query := `SELECT sensor_id, kunstwerk_id, metingen_processed, afwijkingen_detected, last_meting_id, status 
	          FROM analyze_sensor_metingen()`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("fout bij uitvoeren analyse procedure: %w", err)
	}
	defer rows.Close()

	var results []AnalysisResult
	for rows.Next() {
		var result AnalysisResult
		err := rows.Scan(
			&result.SensorID,
			&result.KunstwerkID,
			&result.MetingenProcessed,
			&result.AfwijkingenDetected,
			&result.LastMetingID,
			&result.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("fout bij scannen analyse resultaat: %w", err)
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("fout bij itereren analyse resultaten: %w", err)
	}

	return results, nil
}

// GetErrorLog retrieves recent errors from the procedure error log
func (r *AnalysisProcedureRepository) GetErrorLog(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	query := `SELECT id, logged_at, procedure_name, sensor_id, error_message, error_context 
	          FROM procedure_error_log 
	          ORDER BY logged_at DESC 
	          LIMIT $1`

	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("fout bij ophalen error log: %w", err)
	}
	defer rows.Close()

	var logs []map[string]interface{}
	for rows.Next() {
		var id int64
		var loggedAt interface{}
		var procName string
		var sensorID interface{}
		var errorMsg string
		var errorContext interface{}

		err := rows.Scan(&id, &loggedAt, &procName, &sensorID, &errorMsg, &errorContext)
		if err != nil {
			return nil, fmt.Errorf("fout bij scannen error log: %w", err)
		}

		logs = append(logs, map[string]interface{}{
			"id":              id,
			"logged_at":       loggedAt,
			"procedure_name":  procName,
			"sensor_id":       sensorID,
			"error_message":   errorMsg,
			"error_context":   errorContext,
		})
	}

	return logs, nil
}
