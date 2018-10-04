package server

import (
	"context"
	"time"

	"zephyrus/internal/device"
	"zephyrus/schemas"
)

type WeatherService struct {
	sensor device.Sensor
}

func (s *WeatherService) GetTemperature(ctx context.Context, request *schemas.GetTemperatureRequest) (*schemas.GetTemperatureResponse, error) {
	temperature, err := s.sensor.GetTemperature()
	if err != nil {
		return nil, err
	}

	return &schemas.GetTemperatureResponse{Temperature: temperature}, nil
}

func (s *WeatherService) StreamTemperature(request *schemas.GetTemperatureStreamRequest, stream schemas.Weather_StreamTemperatureServer) error {
	var sample int32

	// The temperature streaming behavior varies based on the number of requested samples:
	//   < 0 -- noop
	//   = 0 -- stream indefinitely
	//   > 0 -- stream only the requested number of samples

	if request.Samples < 0 {
		return nil
	}

	for {
		temperature, err := s.sensor.GetTemperature()
		if err != nil {
			return err
		}

		stream.Send(&schemas.GetTemperatureResponse{Temperature: temperature})

		if request.Samples > 0 {
			sample++

			if request.Samples == sample {
				break
			}
		}

		// Throttle device reads when a sample rate is provided; otherwise, stream readings
		// to the client as fast as it can receive them.
		if request.SampleRate > 0 {
			time.Sleep(time.Duration(1.0e9 / request.SampleRate))
		}
	}

	return nil
}
