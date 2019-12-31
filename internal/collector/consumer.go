package collector

import (
	"fmt"

	"lib.kevinlin.info/aperture"
)

// GlobalMetricNamespace is the top-level namespace for all emitted metrics.
const GlobalMetricNamespace = "zephyrus"

// TemperatureStatsdConsumer is a consumer implementing the client.TemperatureConsumer interface
// for emitting consumed temperatures as metrics to statsd.
type TemperatureStatsdConsumer struct {
	// Backing statsd client.
	client aperture.Statsd
	// Device identifier to attach as a tag to all emitted metrics.
	identifier string
}

// NewTemperatureStatsdConsumer creates a new statsd consumer using the specified device identifier
// and remote statsd address.
func NewTemperatureStatsdConsumer(deviceIdentifier string, addr string) (*TemperatureStatsdConsumer, error) {
	client, err := aperture.NewClient(&aperture.Config{
		Address: addr,
		Prefix:  GlobalMetricNamespace,
	})
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
	tags := map[string]interface{}{
		"device": c.identifier,
	}

	c.client.Gauge(metric, 1000.0*temperature, tags)

	return nil
}
