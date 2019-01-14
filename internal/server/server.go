package server

import (
	"fmt"
	"net"

	"zephyrus/internal/device"
	"zephyrus/schemas"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// ZephyrusServer wraps a gRPC server and registers all necessary services.
type ZephyrusServer struct {
	// Wrapped gRPC server instance.
	server *grpc.Server
}

// NewZephyrusServer creates a new server with the specified device sensor backend.
// Note that the server is, in itself, agnostic to the actual hardware device; it merely provides
// abstractions on top of a client library that implements the device.Sensor interface.
func NewZephyrusServer(sensor device.Sensor) (*ZephyrusServer, error) {
	grpcServer := grpc.NewServer()
	deviceInfoService := &DeviceInfoService{sensor}
	weatherService := &WeatherService{sensor}
	metaService := &MetaService{}

	schemas.RegisterDeviceInfoServer(grpcServer, deviceInfoService)
	schemas.RegisterWeatherServer(grpcServer, weatherService)
	schemas.RegisterMetaServer(grpcServer, metaService)
	reflection.Register(grpcServer)

	return &ZephyrusServer{server: grpcServer}, nil
}

// Serve starts the gRPC server on the specified port and serves indefinitely.
func (s *ZephyrusServer) Serve(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("server: %v", err)
	}

	defer listener.Close()

	if err := s.server.Serve(listener); err != nil {
		return fmt.Errorf("server: %v", err)
	}

	return nil
}
