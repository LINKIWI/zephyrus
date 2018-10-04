package device

import (
	"zephyrus/schemas"
)

// Sensor describes a high-level interface for operations that a hardware temperature sensor
// device might support. This provides a contract for logic higher in the stack to interact with
// device APIs without concerning itself with the nuances of how any specific hardware device
// performs any of these operations.
type Sensor interface {
	// Open opens a connection or communication channel with the device.
	Open() error

	// Close closes a connection or communication channel with the device.
	Close() error

	// GetIdentifier returns a unique string identifier for the device, if available.
	GetIdentifier() (string, error)

	// GetStatus reports the current state of the device.
	GetStatus() schemas.Status

	// GetTemperature returns a live temperature sensor reading from the device.
	GetTemperature() (float64, error)
}
