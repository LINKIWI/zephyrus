package server

import (
	"context"

	"zephyrus/schemas"
)

// MetaService is a server-side implementation of meta RPC methods.
type MetaService struct{}

// HealthCheck always returns a successful health check response.
// Note that this method only reports liveness of the server, without regard to any other components
// of the system (including functionality of attached weather devices).
func (s *MetaService) HealthCheck(ctx context.Context, request *schemas.HealthCheckRequest) (*schemas.HealthCheckResponse, error) {
	return &schemas.HealthCheckResponse{Ok: true}, nil
}
