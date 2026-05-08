package main

import (
	"fmt"
	"log"
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
	consumers.StartConsumingSensorData()
}
