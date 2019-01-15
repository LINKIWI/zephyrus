package device

import (
	"time"

	"zephyrus/internal/cache"
	"zephyrus/schemas"
)

const (
	// Cache key used to identify temperature values from the device.
	temperatureCacheKey = "sensor:temperature"
	// Default TTL for cached temperature values.
	temperatureCacheTTL = 500 * time.Millisecond
)

// ThrottledSensor implements the Sensor interface and wraps another Sensor, throttling request
// volume to the actual device by protecting reads with an in-memory cache.
type ThrottledSensor struct {
	sensor Sensor
	cache  *cache.MemoryTTLCache
}

// NewThrottledSensor creates a throttled sensor from another implementation of the same interface.
func NewThrottledSensor(sensor Sensor) *ThrottledSensor {
	return &ThrottledSensor{
		sensor: sensor,
		cache:  cache.NewMemoryTTLCache(),
	}
}

// Open is proxied directly to the sensor.
func (s *ThrottledSensor) Open() error {
	return s.sensor.Open()
}

// Close is proxied directly to the sensor.
func (s *ThrottledSensor) Close() error {
	return s.sensor.Close()
}

// GetIdentifier is proxied directly to the sensor.
func (s *ThrottledSensor) GetIdentifier() (string, error) {
	return s.sensor.GetIdentifier()
}

// GetStatus is proxied directly to the sensor.
func (s *ThrottledSensor) GetStatus() schemas.Status {
	return s.sensor.GetStatus()
}

// GetTemperature wraps the sensor's equivalent method behind a cache with a default TTL of
// temperatureCacheTTL. This guarantees an upper cap on the QPS made to the actual sensor device.
func (s *ThrottledSensor) GetTemperature() (float64, error) {
	cached := s.cache.Get(temperatureCacheKey)
	if cached != nil {
		return cached.(float64), nil
	}

	temperature, err := s.sensor.GetTemperature()
	defer func() {
		if err == nil {
			s.cache.Set(temperatureCacheKey, temperature, temperatureCacheTTL)
		}
	}()

	return temperature, err
}
