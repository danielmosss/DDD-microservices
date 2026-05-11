package analyse

import (
	"context"
	"log"
	"time"

	"monitoring/internal/infra/db"
)

type AnalysisScheduler struct {
	analysisProcedure *db.AnalysisProcedureRepository
	ticker            *time.Ticker
	stopCh            chan bool
	interval          time.Duration
}

func NewAnalysisScheduler(analysisProcedure *db.AnalysisProcedureRepository, interval time.Duration) *AnalysisScheduler {
	return &AnalysisScheduler{
		analysisProcedure: analysisProcedure,
		interval:          interval,
		stopCh:            make(chan bool),
	}
}

// Start begins the scheduled analysis loop (every interval, default 1 minute)
func (s *AnalysisScheduler) Start(ctx context.Context) {
	s.ticker = time.NewTicker(s.interval)

	go func() {
		defer s.ticker.Stop()

		for {
			select {
			case <-s.stopCh:
				log.Println("Analysis scheduler stopped")
				return
			case <-s.ticker.C:
				s.performAnalysis(ctx)
			}
		}
	}()

	log.Printf("Analysis scheduler started (interval: %v)", s.interval)
}

// Stop halts the scheduled analysis
func (s *AnalysisScheduler) Stop() {
	select {
	case s.stopCh <- true:
	default:
	}
}

func (s *AnalysisScheduler) performAnalysis(ctx context.Context) {
	log.Println("Starting scheduled analysis...")

	results, err := s.analysisProcedure.ExecuteAnalysis(ctx)
	if err != nil {
		log.Printf("Error executing analysis procedure: %v", err)
		return
	}

	log.Printf("Analysis completed: processed %d sensors", len(results))

	for _, result := range results {
		if result.Status != "success" {
			log.Printf("Sensor %d analysis error: %s", result.SensorID, result.Status)
			continue
		}

		if result.AfwijkingenDetected > 0 {
			log.Printf("Sensor %d: %d anomalies detected from %d measurements", result.SensorID, result.AfwijkingenDetected, result.MetingenProcessed)
		}
	}
}

func (s *AnalysisScheduler) GetProcedureErrors(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	return s.analysisProcedure.GetErrorLog(ctx, limit)
}
