package server

import (
	"context"
	"time"

	"zephyrus/internal/device"
	"zephyrus/schemas"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// transientFailureLimit describes the number of times the server is permitted to consecutively
// retry a failed stream transmission before failing the stream entirely.
const transientFailureLimit = 5

// transientClientErrors describes error codes sent by the client during a stream that the
// server may gracefully retry on.
var transientClientErrors = []codes.Code{
	codes.Internal,
	codes.ResourceExhausted,
	codes.Unavailable,
}

// WeatherService is a server-side implementation of weather RPC calls.
type WeatherService struct {
	sensor device.Sensor
}

// GetTemperature reads the current temperature.
func (s *WeatherService) GetTemperature(ctx context.Context, request *schemas.GetTemperatureRequest) (*schemas.GetTemperatureResponse, error) {
	temperature, err := s.sensor.GetTemperature()
	if err != nil {
		return nil, err
	}

	return &schemas.GetTemperatureResponse{Temperature: temperature}, nil
}

// StreamTemperature reads from the sensor multiple times and streams each reading individually back
// to the client. The server-side behavior of this method varies based on the client-supplied
// request parameters.
func (s *WeatherService) StreamTemperature(request *schemas.GetTemperatureStreamRequest, stream schemas.Weather_StreamTemperatureServer) error {
	var sample int32

	// Abstraction to gracefully retry a client stream transmission, up to the maximum number of
	// allowable consecutive failures.
	send := func(response *schemas.GetTemperatureResponse) error {
		var retryWrapper func(int) error

		retryWrapper = func(failures int) error {
			if err := stream.Send(response); err != nil {
				if failures < transientFailureLimit {
					for _, retryErr := range transientClientErrors {
						if status.Code(err) == retryErr {
							time.Sleep(1 * time.Second)
							return retryWrapper(failures + 1)
						}
					}
				}

				return err
			}

			return nil
		}

		return retryWrapper(0)
	}

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

		err = send(&schemas.GetTemperatureResponse{Temperature: temperature})
		if err != nil {
			return err
		}

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
