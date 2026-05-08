package main

import (
	"fmt"
	"log"
	"monitoring/internal/infra/messaging"
	"monitoring/internal/interfaces/consumers"
	"monitoring/server"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello World")

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	server.StartDatabaseConnection()

	messaging.StartDailyHealthSummery()

	// Needs to be the last because i put that quit stuff in there
	consumers.StartConsumingSensorData()
}
