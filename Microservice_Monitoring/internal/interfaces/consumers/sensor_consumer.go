package consumers

import (
	"context"
	"encoding/json"
	"log"
	"monitoring/internal/app/ingest"
	"monitoring/internal/domain/models"
	"monitoring/internal/infra/db"
	"monitoring/server"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

func StartConsumingSensorData() {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}

	durableName := "SensorDataWorker_1"
	retryDelay := 5 * time.Second

	metingRepo := db.NewPostgresMetingRepository(server.GetDBPool())
	ingestService := ingest.NewIngestService(metingRepo)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	for {
		nc, err := nats.Connect(natsURL)
		if err != nil {
			log.Printf("[SENSOR-CONSUMER] Kan niet verbinden met NATS (%v). Nieuwe poging over %s...", err, retryDelay)
			time.Sleep(retryDelay)
			continue
		}

		js, err := nc.JetStream()
		if err != nil {
			log.Printf("[SENSOR-CONSUMER] Kan JetStream niet initialiseren (%v). Nieuwe poging over %s...", err, retryDelay)
			nc.Close()
			time.Sleep(retryDelay)
			continue
		}

		sub, err := js.Subscribe("sensor.data", func(m *nats.Msg) {
			log.Printf("[SENSOR-CONSUMER] Nieuw sensor data bericht ontvangen: %s\n", string(m.Data))

			var incData models.IncMeting
			if err := json.Unmarshal(m.Data, &incData); err != nil {
				log.Printf("[SENSOR-CONSUMER] Fout bij unmarshallen van sensor data: %v\n", err)
				m.Nak()
				return
			}

			ctx := context.Background()
			if _, err := ingestService.VerwerkMeting(ctx, incData); err != nil {
				log.Printf("[SENSOR-CONSUMER] Fout bij verwerken van sensor data: %v", err)
				m.NakWithDelay(10 * time.Second)
				return
			}

			m.Ack()
			log.Printf("[SENSOR-CONSUMER] Sensor data afgehandeld (Ack verzonden).")
		}, nats.Durable(durableName))

		if err != nil {
			log.Printf("[SENSOR-CONSUMER] Fout bij subscriben (%v). Nieuwe poging over %s...", err, retryDelay)
			nc.Close()
			time.Sleep(retryDelay)
			continue
		}

		log.Printf("[SENSOR-CONSUMER] Gestart en luistert op 'sensor.data' (Durable: %s)...\n", durableName)

		<-quit
		_ = sub.Drain()
		_ = nc.Drain()
		return
	}
}
