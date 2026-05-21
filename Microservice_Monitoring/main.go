package main

import (
	"context"
	"fmt"
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
	fmt.Println("Hello World")

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	server.StartDatabaseConnection()

	go messaging.StartDailyHealthSummery()

	analysisRepo := db.NewAnalysisProcedureRepository(server.GetDBPool())
	analysisScheduler := analyse.NewAnalysisScheduler(analysisRepo, time.Minute)
	analysisScheduler.Start(context.Background())

	restapi.StartRestAPI()

	// Needs to be the last because i put that quit stuff in there
	consumers.StartConsumingSensorData()
}
