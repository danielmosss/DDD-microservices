package analyse

import (
	"context"
	"fmt"
	"log"
	"monitoring/internal/db"
	"monitoring/internal/domain/models"
	"monitoring/server"

	"github.com/google/uuid"
)

type ConfiguratieRepository interface {
	GetBySensorID(ctx context.Context, sensorID uuid.UUID) (models.SensorConfiguratie, error)
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
		fmt.Println("Handmatige meting, sla automatische config-check over.")
		return nil
	}

	config, err := s.configRepo.GetBySensorID(ctx, *m.SensorID)
	if err != nil {
		return fmt.Errorf("kan bedrijfsregels niet ophalen voor sensor %s: %w", m.SensorID, err)
	}

	if config.MaxValue != nil && m.Waarde > *config.MaxValue {
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

		log.Printf("%s", afwijking)

		// TODO: Sla de afwijking op in de database
		// TODO: Publiceer het AfwijkingGedetecteerd event naar NATS

		if isWarning {
			fmt.Printf("⚠️ WAARSCHUWING: Meting %.2f valt binnen de %v%% marge van norm %.2f\n", m.Waarde, *config.MargePercentage, *config.MaxValue)
		} else {
			fmt.Printf("🚨 KRITIEK: Meting %.2f overschrijdt de marge van norm %.2f keihard!\n", m.Waarde, *config.MaxValue)
		}
	} else {
		fmt.Printf("✅ Meting %.2f is helemaal in orde.\n", m.Waarde)
	}

	return nil
}
