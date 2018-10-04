package client

// TemperatureConsumer describes a type that asynchronously consumes temperature readings. The
// producer is the gRPC client, via the gRPC server's temperature streaming API.
type TemperatureConsumer interface {
	// Consume a single temperature value, in celsius units.
	// The consumer may optionally return a non-nil error to abort the streaming operation.
	Consume(temperature float64) error
}
