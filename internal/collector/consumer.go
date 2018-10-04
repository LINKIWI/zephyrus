package collector

import (
	"fmt"

	"github.com/cactus/go-statsd-client/statsd"
)

// GlobalMetricNamespace is the top-level namespace for all emitted metrics.
const GlobalMetricNamespace = "zephyrus"

// TemperatureStatsdConsumer is a consumer implementing the client.TemperatureConsumer interface
// for emitting consumed temperatures as metrics to statsd.
type TemperatureStatsdConsumer struct {
	// Backing statsd client.
	client statsd.Statter
	// Device identifier to attach as a tag to all emitted metrics.
	identifier string
}

// NewTemperatureStatsdConsumer creates a new statsd consumer using the specified device identifier
// and remote statsd address.
func NewTemperatureStatsdConsumer(deviceIdentifier string, addr string) (*TemperatureStatsdConsumer, error) {
	client, err := statsd.NewClient(addr, GlobalMetricNamespace)
	if err != nil {
		return nil, fmt.Errorf("consumer: %v", err)
	}

	return &TemperatureStatsdConsumer{
		client:     client,
		identifier: deviceIdentifier,
	}, nil
}

// Consume ships the passed temperature to statsd as a gauge with properly formatted names and tags.
func (c *TemperatureStatsdConsumer) Consume(temperature float64) error {
	metric := "collector.temperature"
	tags := map[string]string{
		"device": c.identifier,
	}

	c.client.Gauge(formatMetric(metric, tags), int64(1000.0*temperature), 1.0)

	return nil
}
