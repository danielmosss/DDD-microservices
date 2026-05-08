package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

var NatsStream nats.JetStreamContext
var natsConn *nats.Conn

func main() {
	StartMessageBroker()
	setupStream()
	StartEventSubscriber()
	StartGeneratingSensorDataAndPublishing()
}

func StartMessageBroker() {
	natsConn, _ = nats.Connect(nats.DefaultURL)

	var jsErr error
	NatsStream, jsErr = natsConn.JetStream()
	if jsErr != nil {
		log.Fatalf("Kan JetStream niet initialiseren: %v", jsErr)
	}
}

func StartGeneratingSensorDataAndPublishing() {
	// Simuleer het genereren van sensor data en publiceer deze elke 5 seconden
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Genereer een willekeurige meting
			randomNumberOneOrTwo := time.Now().Unix()%2 + 1

			meting := Meting{}

			if randomNumberOneOrTwo == 0 {
				meting = Meting{
					SensorID:    1,
					KunstwerkID: 1,
					Waarde:      11.1,
				}
			} else {
				meting = Meting{
					SensorID:    2,
					KunstwerkID: 1,
					Waarde:      29.2,
				}
			}

			// Publiceer de meting op NATS
			PublishMessage("sensor.data", meting)
		}
	}
}

func setupStream() {
	streamName := "SensorData"

	_, err := NatsStream.StreamInfo(streamName)
	if err != nil {
		log.Printf("Stream %s maken...", streamName)
		_, err = NatsStream.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: []string{"sensor.*"},
		})
		if err != nil {
			log.Fatalf("Fout bij aanmaken stream: %v", err)
		}
	} else {
		log.Printf("Stream %s bestaat al.", streamName)
	}
}

func PublishMessage(subject string, data Meting) {
	if NatsStream == nil {
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Fout bij serialiseren van data: %v", err)
	}

	_, err = NatsStream.Publish(subject, payloadBytes)
	if err != nil {
		log.Fatalf("Fout bij publiceren van bericht: %v", err)
	}

	log.Printf("Bericht gepubliceerd op subject %s: %s\n", subject, string(payloadBytes))
}

func StartEventSubscriber() {
	durableName := "NotificationWorker_1"

	_, err := NatsStream.Subscribe(">", func(m *nats.Msg) {
		currentTime := time.Now()
		fmt.Printf("\n%s - [NOTIFICATION-SERVICE] Bericht binnen: %s op subject %s\n", currentTime, string(m.Data), m.Subject)

		m.Ack()
		fmt.Println("[NOTIFICATION-SERVICE] E-mail verstuurd (Ack verzonden).")

	}, nats.Durable(durableName))

	if err != nil {
		log.Fatalf("Fout bij subscriben: %v", err)
	}

	fmt.Printf("[NOTIFICATION-SERVICE] Gestart en luistert op '>' (Durable: %s)\n", durableName)
}
