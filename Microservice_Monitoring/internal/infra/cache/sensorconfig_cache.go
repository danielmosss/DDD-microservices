package cache

import (
	"context"
	"sync"
	"time"

	"monitoring/internal/app/analyse"
	"monitoring/internal/domain/models"
)

// cacheItem houdt de data én de houdbaarheidsdatum vast
type cacheItem struct {
	config  models.SensorConfiguratie
	expires time.Time
}

type CachedConfiguratieRepository struct {
	// De "echte" database repo die we aanroepen bij een cache miss
	next analyse.ConfiguratieRepository

	// Ons in-memory geheugen
	data map[int64]cacheItem
	mu   sync.RWMutex // Voorkomt crashes als NATS met 100 threads tegelijk leest/schrijft
	ttl  time.Duration
}

// Constructor
func NewCachedConfiguratieRepository(next analyse.ConfiguratieRepository, ttl time.Duration) *CachedConfiguratieRepository {
	return &CachedConfiguratieRepository{
		next: next,
		data: make(map[int64]cacheItem),
		ttl:  ttl,
	}
}

// GetBySensorID implementeert de interface
func (c *CachedConfiguratieRepository) GetBySensorID(ctx context.Context, sensorID int64) (models.SensorConfiguratie, error) {
	// 1. Probeer te lezen uit de cache (Read lock is snel en staat meerdere lezers toe)
	c.mu.RLock()
	item, found := c.data[sensorID]
	c.mu.RUnlock()

	// 2. Cache Hit: Is het gevonden én nog niet over de datum?
	if found && time.Now().Before(item.expires) {
		return item.config, nil // Direct teruggeven zonder database call!
	}

	// 3. Cache Miss (of verlopen): Haal op uit de échte database
	config, err := c.next.GetBySensorID(ctx, sensorID)
	if err != nil {
		return models.SensorConfiguratie{}, err
	}

	// 4. Opslaan in de cache voor de volgende keer (Write lock)
	c.mu.Lock()
	c.data[sensorID] = cacheItem{
		config:  config,
		expires: time.Now().Add(c.ttl),
	}
	c.mu.Unlock()

	return config, nil
}
