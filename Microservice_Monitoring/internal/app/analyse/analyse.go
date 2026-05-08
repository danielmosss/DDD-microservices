package analyse

import (
	"context"
	"fmt"
	"log"
	"monitoring/internal/db"
	"monitoring/internal/domain/models"
	"monitoring/server"
)

type ConfiguratieRepository interface {
	GetBySensorID(ctx context.Context, sensorID int64) (models.SensorConfiguratie, error)
}

func AnalyzeIncommingSensorData(meting models.Meting) {
	serverConfig := db.NewPostgresConfiguratieRepository(server.GetDBPool())
	analyzeService := NewAnalyseService(serverConfig)

	ctx := context.Background()
	err := analyzeService.AnalyseerMeting(ctx, meting)

	if err != nil {
		log.Printf("Error in analyse: %v", err)
	}
}

type AnalyseService struct {
	configRepo ConfiguratieRepository
}

func NewAnalyseService(repo ConfiguratieRepository) *AnalyseService {
	return &AnalyseService{
		configRepo: repo,
	}
}

func (s *AnalyseService) AnalyseerMeting(ctx context.Context, m models.Meting) error {
	if m.SensorID == nil {
		return fmt.Errorf("sensor id is nil in meting: %v", m)
	}

	sensorID := *m.SensorID
	config, err := s.configRepo.GetBySensorID(ctx, sensorID)
	if err != nil {
		return fmt.Errorf("kan bedrijfsregels niet ophalen voor sensor %d: %w", sensorID, err)
	}

	if config.MaxValue != nil && m.Waarde > *config.MaxValue {
		if config.MargePercentage == nil {
			config.MargePercentage = new(float64)
			*config.MargePercentage = 0
		}

		margeAbsoluut := (*config.MaxValue / 100.0) * *config.MargePercentage
		grensWaarschuwing := *config.MaxValue + margeAbsoluut

		isWarning := m.Waarde <= grensWaarschuwing

		afwijking := models.Afwijking{
			MetingID:      m.ID,
			KunstwerkID:   m.KunstwerkID,
			Time:          m.Time,
			NormWaarde:    *config.MaxValue,
			GemetenWaarde: m.Waarde,
			IsWarning:     isWarning,
		}

		log.Printf("Afwijking gedetecteerd: %+v", afwijking)

		// TODO: Sla de afwijking op in de database
		// TODO: Publiceer het AfwijkingGedetecteerd event naar NATS

		if isWarning {
			log.Printf("⚠️ WAARSCHUWING: Meting %.2f valt binnen de %.2f%% marge van norm %.2f", m.Waarde, *config.MargePercentage, *config.MaxValue)
		} else {
			log.Printf("🚨 KRITIEK: Meting %.2f overschrijdt de marge van norm %.2f keihard!", m.Waarde, *config.MaxValue)
		}
	} else {
		log.Printf("✅ Meting %.2f is helemaal in orde.", m.Waarde)
	}

	return nil
}
