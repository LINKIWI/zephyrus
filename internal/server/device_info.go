package server

import (
	"context"

	"zephyrus/internal/device"
	"zephyrus/schemas"
)

type DeviceInfoService struct {
	sensor device.Sensor
}

func (s *DeviceInfoService) GetIdentifier(ctx context.Context, request *schemas.GetIdentifierRequest) (*schemas.GetIdentifierResponse, error) {
	identifier, err := s.sensor.GetIdentifier()
	if err != nil {
		return nil, err
	}

	return &schemas.GetIdentifierResponse{Identifier: identifier}, nil
}

func (s *DeviceInfoService) GetStatus(ctx context.Context, request *schemas.GetStatusRequest) (*schemas.GetStatusResponse, error) {
	status := s.sensor.GetStatus()

	return &schemas.GetStatusResponse{Status: status}, nil
}
