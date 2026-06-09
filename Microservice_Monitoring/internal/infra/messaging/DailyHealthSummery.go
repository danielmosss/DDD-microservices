package messaging

import (
	"context"
	"log"
	"monitoring/internal/domain/events"
	"monitoring/internal/infra/db"
	"monitoring/server"
	"time"
)

// create a function that will start the DailyHealthSummary it will send every 24 hours a report about a kunstwerk.
// it will check every 5 minuts in the database which still need to be done with a limit of 5 per check. For now i will just print the report to the console but in the future it will be send to a message queue.
func StartDailyHealthSummary() {
	// do one immediate check on startup, then every 5 minutes
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	checkAndGenerate := func() {
		// check the database for kunstwerken that need a report
		kunstwerken, err := getKunstwerkenNeedingReport()
		if err != nil {
			log.Printf("Error fetching kunstwerken: %v", err)
			return
		}

		for _, kunstwerkId := range kunstwerken {
			generateReport(kunstwerkId)
		}
	}

	// immediate run on startup
	checkAndGenerate()

	for {
		select {
		case <-ticker.C:
			checkAndGenerate()
		}
	}
}

func getKunstwerkenNeedingReport() ([]int64, error) {
	KunstwerkPostgres := db.NewPostgresKunstwerkRepository(server.GetDBPool())
	ctx := context.Background()
	kunstwerken, err := KunstwerkPostgres.GetKunstwerkenNeedingReport(ctx)
	if err != nil {
		return nil, err
	}

	return kunstwerken, nil
}

func generateReport(kunstwerkId int64) {
	KunstwerkPostgres := db.NewPostgresKunstwerkRepository(server.GetDBPool())
	ctx := context.Background()

	DHS, err := events.NewDailyHealthSummary(kunstwerkId)
	if err != nil {
		log.Printf("Error generating Daily Health Summary for Kunstwerk ID %d: %v", kunstwerkId, err)
		return
	}
	log.Printf("Daily Health Summary for Kunstwerk ID %d: %+v", kunstwerkId, DHS)

	err = KunstwerkPostgres.UpdateKunstwerkDHupdateTime(ctx, kunstwerkId)
	if err != nil {
		log.Printf("Error updating last_send_dh_update for Kunstwerk ID %d: %v", kunstwerkId, err)
	}

	// TODO: later implement sending report as an event to the message queue instead of printing it to the console
}
