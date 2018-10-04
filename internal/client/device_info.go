package client

import (
	"context"
	"fmt"

	"zephyrus/schemas"
)

type DeviceInfoService struct {
	client schemas.DeviceInfoClient
}

func (s *DeviceInfoService) GetIdentifier() (string, error) {
	ctx := context.Background()
	req := &schemas.GetIdentifierRequest{}

	resp, err := s.client.GetIdentifier(ctx, req)
	if err != nil {
		return "", fmt.Errorf("device_info: %v", err)
	}

	return resp.Identifier, nil
}

func (s *DeviceInfoService) GetStatus() (schemas.Status, error) {
	ctx := context.Background()
	req := &schemas.GetStatusRequest{}

	resp, err := s.client.GetStatus(ctx, req)
	if err != nil {
		return schemas.Status_UNKNOWN, fmt.Errorf("device_info: %v", err)
	}

	return resp.Status, nil
}
