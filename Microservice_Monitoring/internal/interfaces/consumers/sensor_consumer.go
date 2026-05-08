package consumers

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	_, err = js.Subscribe("sensor.data", func(m *nats.Msg) {
		log.Printf("[SENSOR-CONSUMER] Nieuw sensor data bericht ontvangen: %s\n", string(m.Data))

		//inc = IncMeting struct
		var IncData IncMeting
		err = json.Unmarshal(m.Data, &IncData)
		if err != nil {
			log.Printf("[SENSOR-CONSUMER] Fout bij unmarshallen van sensor data: %v\n", err)
			m.Nak() // Vertel NATS: het is mislukt, stuur later opnieuw
			return
		}

		//log for debug
		log.Printf("KunstwerkID: %s\n", IncData.KunstwerkID)
		log.Printf("SensorID: %s\n", IncData.SensorID)
		log.Printf("Waarde: %s\n", IncData.Waarde)

		// Hier verwerk je de sensor data. Bijvoorbeeld:
		// err := processSensorData(m.Data)
		// if err != nil {
		//     m.Nak() // Vertel NATS: het is mislukt, stuur later opnieuw
		//     return
		// }

		m.Ack()
		log.Printf("[SENSOR-CONSUMER] Sensor data netjes afgehandeld (Ack verzonden).\n")

	}, nats.Durable(durableName))

	if err != nil {
		log.Fatalf("Fout bij subscriben: %v", err)
	}

	log.Printf("[SENSOR-CONSUMER] Gestart en luistert op 'sensor.data' (Durable: %s)...\n", durableName)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
