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
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Kan niet verbinden met NATS: %v", err)
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("Kan JetStream niet initialiseren: %v", err)
	}

	durableName := "SensorDataWorker_1"

	metingRepo := db.NewPostgresMetingRepository(server.GetDBPool())
	ingestService := ingest.NewIngestService(metingRepo)

	_, err = js.Subscribe("sensor.data", func(m *nats.Msg) {
		log.Printf("[SENSOR-CONSUMER] Nieuw sensor data bericht ontvangen: %s\n", string(m.Data))

		//inc = IncMeting struct
		var IncData models.IncMeting
		err = json.Unmarshal(m.Data, &IncData)
		if err != nil {
			log.Printf("[SENSOR-CONSUMER] Fout bij unmarshallen van sensor data: %v\n", err)
			m.Nak() // Vertel NATS: het is mislukt, stuur later opnieuw
			return
		}

		ctx := context.Background()
		_, err := ingestService.VerwerkMeting(ctx, IncData)

		// log for debug (use correct verbs and handle nil SensorID)
		log.Printf("KunstwerkID: %d", IncData.KunstwerkID)
		if IncData.SensorID != nil {
			log.Printf("SensorID: %d", *IncData.SensorID)
		} else {
			log.Printf("SensorID: <nil>")
		}
		log.Printf("Waarde: %.6f", IncData.Waarde)

		if err != nil {
			log.Printf("[SENSOR-CONSUMER] Fout bij verwerken van sensor data: %v", err)
			m.NakWithDelay(10 * time.Second)
		} else {
			m.Ack()
			log.Printf("[SENSOR-CONSUMER] Sensor data netjes afgehandeld (Ack verzonden).")
		}
	}, nats.Durable(durableName))

	if err != nil {
		log.Fatalf("Fout bij subscriben: %v", err)
	}

	log.Printf("[SENSOR-CONSUMER] Gestart en luistert op 'sensor.data' (Durable: %s)...\n", durableName)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
