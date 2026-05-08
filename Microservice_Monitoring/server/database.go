package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	dbPool *pgxpool.Pool
	once   sync.Once
)

func StartDatabaseConnection() {
	once.Do(func() {
		dbURL := os.Getenv("DATABASE_URL")
		if dbURL == "" {
			log.Fatal("DATABASE_URL environment variable is niet gezet")
		}

		ctx := context.Background()
		pool, err := pgxpool.New(ctx, dbURL)
		if err != nil {
			log.Fatalf("Kan geen verbinding maken met de database: %v", err)
		}

		if err := pool.Ping(ctx); err != nil {
			pool.Close()
			log.Fatalf("Database ping mislukt: %v", err)
		}

		dbPool = pool
		fmt.Println("Database connectie succesvol!")
	})
}

func GetDBPool() *pgxpool.Pool {
	return dbPool
}

func CloseDBPool() {
	if dbPool != nil {
		dbPool.Close()
		dbPool = nil
	}
}
