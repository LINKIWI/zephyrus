package server

import (
	"context"

	"zephyrus/internal/device"
	"zephyrus/schemas"
)

// DeviceInfoService is a server-side implementation of device information RPC calls.
type DeviceInfoService struct {
	sensor device.Sensor
}

// GetIdentifier gets the device identifier.
func (s *DeviceInfoService) GetIdentifier(ctx context.Context, request *schemas.GetIdentifierRequest) (*schemas.GetIdentifierResponse, error) {
	identifier, err := s.sensor.GetIdentifier()
	if err != nil {
		return nil, err
	}

	return &schemas.GetIdentifierResponse{Identifier: identifier}, nil
}

// GetStatus gets the current status of the device.
func (s *DeviceInfoService) GetStatus(ctx context.Context, request *schemas.GetStatusRequest) (*schemas.GetStatusResponse, error) {
	status := s.sensor.GetStatus()

	return &schemas.GetStatusResponse{Status: status}, nil
}
