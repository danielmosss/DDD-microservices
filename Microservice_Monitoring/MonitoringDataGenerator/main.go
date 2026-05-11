package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
)

var NatsStream nats.JetStreamContext
var natsConn *nats.Conn
var dbPool *pgxpool.Pool

type Sensor struct {
	ID          int64
	SensorID    int64
	TypeID      int32
	MinValue    *float64
	MaxValue    *float64
	Margin      *float64
	KunstwerkID int64
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	// Initialize database connection
	var err error

	// Connection string - adjust based on where you're running this
	// From local machine: localhost:5432
	// From Docker container: timescaledb:5432
	connString := "postgres://user:password@localhost:5432/monitoring?sslmode=disable"

	dbPool, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("Fout bij database connectie: %v", err)
	}
	defer dbPool.Close()

	// Test connection
	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatalf("Database niet beschikbaar: %v", err)
	}
	log.Println("Database verbinding OK")

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
	log.Println("Laden van sensoren en configuraties uit database...")

	sensors, err := LoadSensorsWithConfig(context.Background())
	if err != nil {
		log.Fatalf("Fout bij laden sensoren: %v", err)
	}

	if len(sensors) == 0 {
		log.Fatalf("Geen sensoren gevonden in database")
	}

	log.Printf("Geladen: %d sensoren met configuraties", len(sensors))

	// Ticker for 200ms intervals
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	sensorIndex := 0

	for {
		select {
		case <-ticker.C:
			// Publish one sensor per tick (round-robin)
			sensor := sensors[sensorIndex%len(sensors)]

			// Generate realistic data based on sensor type and config
			value := GenerateRealisticSensorValue(sensor)

			meting := Meting{
				SensorID:    sensor.SensorID,
				KunstwerkID: sensor.KunstwerkID,
				Waarde:      value,
			}

			PublishMessage("sensor.data", meting)
			sensorIndex++
		}
	}
}

// LoadSensorsWithConfig queries database for all active sensors with their configurations
func LoadSensorsWithConfig(ctx context.Context) ([]Sensor, error) {
	query := `
		SELECT 
			s.id,
			s.sensortype_id,
			s.kunstwerk_id,
			COALESCE(sc.min_value, 0) as min_value,
			COALESCE(sc.max_value, 0) as max_value,
			COALESCE(sc.marge_percentage, 0) as marge_percentage
		FROM sensor s
		LEFT JOIN sensorconfiguratie sc ON s.id = sc.sensor_id
		WHERE s.deleted = false
		ORDER BY s.id
	`

	rows, err := dbPool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sensors []Sensor
	for rows.Next() {
		var s Sensor
		var minVal, maxVal, marge float64

		err := rows.Scan(&s.SensorID, &s.TypeID, &s.KunstwerkID, &minVal, &maxVal, &marge)
		if err != nil {
			return nil, err
		}

		s.MinValue = &minVal
		s.MaxValue = &maxVal
		s.Margin = &marge
		sensors = append(sensors, s)
	}

	return sensors, rows.Err()
}

// GenerateRealisticSensorValue generates realistic measurement values with occasional anomalies
func GenerateRealisticSensorValue(sensor Sensor) float64 {
	if sensor.MaxValue == nil || *sensor.MaxValue == 0 {
		// Single threshold sensor - generate around min_value
		baseValue := *sensor.MinValue
		variation := (rand.Float64() - 0.5) * baseValue * 0.2 // ±10% variation

		// 15% chance of anomaly
		if rand.Float64() < 0.15 {
			// Generate anomaly: significantly exceed threshold
			anomaly := baseValue * (1.5 + rand.Float64()*2)
			return baseValue + anomaly
		}

		return baseValue + variation
	}

	// Range-based sensor
	minVal := *sensor.MinValue
	maxVal := *sensor.MaxValue
	midpoint := (minVal + maxVal) / 2
	halfRange := (maxVal - minVal) / 2

	// Generate normal distribution around midpoint
	stdDev := halfRange / 3
	normalValue := midpoint + (randGaussian() * stdDev)

	// Clamp to reasonable range
	normalValue = math.Max(normalValue, minVal-halfRange*0.3)
	normalValue = math.Min(normalValue, maxVal+halfRange*0.3)

	// 10% chance of anomaly (outside normal range)
	if rand.Float64() < 0.10 {
		if rand.Float64() < 0.5 {
			// Too low anomaly
			return minVal - halfRange*0.5
		}
		// Too high anomaly
		return maxVal + halfRange*0.5
	}

	return normalValue
}

// randGaussian generates a random number from normal distribution (Box-Muller transform)
func randGaussian() float64 {
	u1 := rand.Float64()
	u2 := rand.Float64()
	z0 := math.Sqrt(-2.0*math.Log(u1)) * math.Cos(2.0*math.Pi*u2)
	return z0
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
		fmt.Printf("%s - [NOTIFICATION-SERVICE] Bericht binnen: %s op subject %s\n", currentTime, string(m.Data), m.Subject)
		m.Ack()
	}, nats.Durable(durableName))

	if err != nil {
		log.Fatalf("Fout bij subscriben: %v", err)
	}

	fmt.Printf("[NOTIFICATION-SERVICE] Gestart en luistert op '>' (Durable: %s)\n", durableName)
}
