package main

import (
	"context"
	"log"
	"monitoring/internal/app/analyse"
	"monitoring/internal/app/restapi"
	"monitoring/internal/infra/db"
	"monitoring/internal/infra/messaging"
	"monitoring/internal/interfaces/consumers"
	"monitoring/server"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf(".env file not loaded (%v); using environment variables", err)
	}

	server.StartDatabaseConnection()

	go messaging.StartDailyHealthSummary()

	analysisRepo := db.NewAnalysisProcedureRepository(server.GetDBPool())
	analysisScheduler := analyse.NewAnalysisScheduler(analysisRepo, time.Second*10)
	analysisScheduler.Start(context.Background())

	go restapi.StartRestAPI()

	// Needs to be the last because i put that quit stuff in there
	consumers.StartConsumingSensorData()
}
