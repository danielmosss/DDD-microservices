package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
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
	// 1. Maak connectie met je NATS broker
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Kan niet verbinden met NATS: %v", err)
	}
	natsConn = nc

	// 2. Initialiseer de JetStream context
	var jsErr error
	NatsStream, jsErr = nc.JetStream()
	if jsErr != nil {
		log.Fatalf("Kan JetStream niet initialiseren: %v", jsErr)
	}
}

func StartGeneratingSensorDataAndPublishing() {
	// Simuleer het genereren van sensor data en publiceer deze elke 5 seconden
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Genereer een willekeurige meting (hier hardcoded voor demo)
			meting := Meting{
				SensorID:    "33333333-33333333-3333-333333333333",
				KunstwerkID: "11111111-1111-1111-1111-111111111111",
				Waarde:      int(rand.Float64() * 30),
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
