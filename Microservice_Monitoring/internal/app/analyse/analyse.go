package analyse

import (
	"context"
	"fmt"
	"log"
	"monitoring/internal/db"
	"monitoring/internal/domain/models"
	"monitoring/server"
	"time"
)

type Status string

const (
	StatusOke     Status = "oke"
	StatusWarning Status = "warning"
	StatusFatal   Status = "fatal"
)

type AdjustWaarde string

const (
	Positief AdjustWaarde = "positief"
	Negatief AdjustWaarde = "negatief"
)

type ConfiguratieRepository interface {
	GetBySensorID(ctx context.Context, sensorID int64) (models.SensorConfiguratie, error)
}

type AfwijkingRepository interface {
	Save(ctx context.Context, m models.Afwijking) (models.Afwijking, error)
}

func AnalyzeIncommingSensorData(meting models.Meting) error {
	serverConfig := db.NewPostgresConfiguratieRepository(server.GetDBPool())
	afwijkingRepo := db.NewPostgresAfwijkingRepository(server.GetDBPool())
	analyzeService := NewAnalyseService(serverConfig, afwijkingRepo)

	ctx := context.Background()
	err := analyzeService.AnalyseerMeting(ctx, meting)

	if err != nil {
		return fmt.Errorf("analyse meting: %v", err)
	}
	return nil
}

type AnalyseService struct {
	configRepo    ConfiguratieRepository
	afwijkingRepo AfwijkingRepository
}

func NewAnalyseService(repo ConfiguratieRepository, afwijkingRepo AfwijkingRepository) *AnalyseService {
	return &AnalyseService{
		configRepo:    repo,
		afwijkingRepo: afwijkingRepo,
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

	// S1 - Simpel eenvoudige drempelwaarde, min_value met of zonder marge_percentage
	if config.MaxValue == nil || *config.MaxValue == 0 {
		outputStatus := CheckDempelWaarde(m.Waarde, config)
		if outputStatus != StatusOke {
			isWarning := outputStatus == StatusWarning

			afwijking := models.Afwijking{
				MetingID:      m.ID,
				MetingTime:    m.Time,
				KunstwerkID:   m.KunstwerkID,
				Time:          time.Now(),
				NormMinWaarde: *config.MinValue,
				NormMaxWaarde: 0,
				GemetenWaarde: m.Waarde,
				IsWarning:     isWarning,
			}

			saved, err := s.afwijkingRepo.Save(ctx, afwijking)
			if err != nil {
				return fmt.Errorf("fout bij opslaan afwijking: %w", err)
			}
			log.Printf("Afwijking opgeslagen: %+v", saved)
		}
	} else {
		// S2 - Wanneer range drempelwaarde. met of zonder marge_percentage
		outputStatus := CheckDrempelRange(m.Waarde, config)
		if outputStatus != StatusOke {
			isWarning := outputStatus == StatusWarning

			afwijking := models.Afwijking{
				MetingID:      m.ID,
				MetingTime:    m.Time,
				KunstwerkID:   m.KunstwerkID,
				Time:          time.Now(),
				NormMinWaarde: *config.MinValue,
				NormMaxWaarde: *config.MaxValue,
				GemetenWaarde: m.Waarde,
				IsWarning:     isWarning,
			}

			saved, err := s.afwijkingRepo.Save(ctx, afwijking)
			if err != nil {
				return fmt.Errorf("fout bij opslaan afwijking: %w", err)
			}
			log.Printf("Afwijking opgeslagen: %+v", saved)
		}
	}

	return nil
}

func CheckDempelWaarde(waarde float64, config models.SensorConfiguratie) Status {
	if config.MargePercentage == nil || *config.MargePercentage == 0 {
		if waarde > *config.MinValue {
			return StatusFatal
		}
		if waarde < *config.MinValue {
			return StatusFatal
		}
	} else {
		minValueAdjusted := valueAdjust(*config.MinValue, config, Positief)
		minValueNegAdjusted := valueAdjust(*config.MinValue, config, Negatief)
		if waarde > minValueAdjusted || waarde < minValueNegAdjusted {
			return StatusFatal
		}
		if waarde < minValueAdjusted && waarde > *config.MinValue {
			return StatusWarning
		}

		if waarde > minValueNegAdjusted && waarde < *config.MinValue {
			return StatusWarning
		}
	}

	return StatusOke
}

func CheckDrempelRange(waarde float64, config models.SensorConfiguratie) Status {
	// Als binnen het bereik dan direct status oke.
	if waarde >= *config.MinValue && waarde <= *config.MaxValue {
		return StatusOke
	}

	//Nu is dus buitenberijk, kijken of fatal of warning
	if config.MargePercentage == nil || *config.MargePercentage == 0 {
		if waarde > *config.MaxValue || waarde < *config.MinValue {
			return StatusFatal
		}
	} else {
		minValueAdjusted := valueAdjust(*config.MinValue, config, Negatief)
		maxValueAdjusted := valueAdjust(*config.MaxValue, config, Positief)

		if waarde >= minValueAdjusted && waarde <= *config.MinValue {
			return StatusWarning
		}

		if waarde <= maxValueAdjusted && waarde >= *config.MaxValue {
			return StatusWarning
		}

		return StatusFatal
	}
	return StatusOke
}

func valueAdjust(waarde float64, config models.SensorConfiguratie, adjustType AdjustWaarde) float64 {
	waardeAdjusted := waarde

	if adjustType == Positief {
		waardeAdjusted = waarde + (waarde * *config.MargePercentage / 100.0)
	}
	if adjustType == Negatief {
		waardeAdjusted = waarde - (waarde * *config.MargePercentage / 100.0)
	}
	return waardeAdjusted
}
