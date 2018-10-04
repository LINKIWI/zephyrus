package client

import (
	"fmt"

	"zephyrus/schemas"

	"google.golang.org/grpc"
)

// ZephyrusClient is an abstraction over a gRPC client to the Zephyrus gRPC server.
type ZephyrusClient struct {
	// DeviceInfo provides abstractions for reading device metadata.
	DeviceInfo *DeviceInfoService
	// Weather provides abstractions for reading and streaming weather information.
	Weather *WeatherService

	// The underlying persistent connection to the server.
	conn *grpc.ClientConn
}

// NewZephyrusClient creates a new client instance for the server at the specified address.
func NewZephyrusClient(addr string) (*ZephyrusClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("client: %v", err)
	}

	deviceInfo := &DeviceInfoService{client: schemas.NewDeviceInfoClient(conn)}
	weather := &WeatherService{client: schemas.NewWeatherClient(conn)}

	return &ZephyrusClient{
		DeviceInfo: deviceInfo,
		Weather:    weather,
		conn:       conn,
	}, nil
}

// Close closes the active connection.
func (c *ZephyrusClient) Close() error {
	return c.conn.Close()
}
