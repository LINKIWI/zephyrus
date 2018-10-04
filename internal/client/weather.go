package client

import (
	"context"
	"fmt"
	"io"

	"zephyrus/schemas"
)

// WeatherService provides client abstractions over the weather service.
type WeatherService struct {
	client schemas.WeatherClient
}

// GetTemperature reads the current temperature.
func (s *WeatherService) GetTemperature() (float64, error) {
	ctx := context.Background()
	req := &schemas.GetTemperatureRequest{}

	resp, err := s.client.GetTemperature(ctx, req)
	if err != nil {
		return 0.0, fmt.Errorf("weather: %v", err)
	}

	return resp.Temperature, nil
}

// StreamTemperature continuously and indefinitely streams temperature readings at a specified
// server-side sample rate.
func (s *WeatherService) StreamTemperature(sampleRate float64, consumer TemperatureConsumer) error {
	return s.StreamTemperatureSamples(sampleRate, 0, consumer)
}

// StreamTemperatureSamples requests a stream of a specified number of samples at s specified sample
// rate.
func (s *WeatherService) StreamTemperatureSamples(sampleRate float64, samples int32, consumer TemperatureConsumer) error {
	ctx := context.Background()
	req := &schemas.GetTemperatureStreamRequest{
		Samples:    samples,
		SampleRate: sampleRate,
	}

	stream, err := s.client.StreamTemperature(ctx, req)
	if err != nil {
		return fmt.Errorf("weather: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("weather: %v", err)
		}

		if err := consumer.Consume(resp.Temperature); err != nil {
			return fmt.Errorf("weather: %v", err)
		}
	}

	return nil
}
