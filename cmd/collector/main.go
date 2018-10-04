package main

import (
	"errors"
	"flag"
	"log"

	"zephyrus/internal/client"
	"zephyrus/internal/collector"
)

type config struct {
	ZephyrusAddr string
	StatsdAddr   string
	SampleRate   float64
}

func main() {
	cfg, err := parseConfig()
	if err != nil {
		panic(err)
	}

	log.Printf(
		"collector: using configuration: zephyrus=%s statsd=%s sample rate=%f",
		cfg.ZephyrusAddr,
		cfg.StatsdAddr,
		cfg.SampleRate,
	)

	log.Printf("collector: connecting to Zephyrus gRPC server")
	zephyrus, err := client.NewZephyrusClient(cfg.ZephyrusAddr)
	if err != nil {
		panic(err)
	}
	defer zephyrus.Close()

	log.Printf("collector: reading device metadata")
	identifier, err := zephyrus.DeviceInfo.GetIdentifier()
	if err != nil {
		panic(err)
	}

	log.Printf("collector: connecting to statsd server")
	consumer, err := collector.NewTemperatureStatsdConsumer(identifier, cfg.StatsdAddr)
	if err != nil {
		panic(err)
	}

	log.Printf("collector: starting collection")
	if err := zephyrus.Weather.StreamTemperature(cfg.SampleRate, consumer); err != nil {
		panic(err)
	}
}

func parseConfig() (*config, error) {
	zephyrusAddr := flag.String("zephyrus", "", "Address of the Zephyrus gRPC server")
	statsdAddr := flag.String("statsd", "", "Address of the statsd server")
	sampleRate := flag.Float64("sample-rate", 1.0, "Collection sample rate from the device")
	flag.Parse()

	if *zephyrusAddr == "" {
		return nil, errors.New("config: address of Zephyrus server must be specified")
	}

	if *statsdAddr == "" {
		return nil, errors.New("config: address of statsd server must be specified")
	}

	return &config{
		ZephyrusAddr: *zephyrusAddr,
		StatsdAddr:   *statsdAddr,
		SampleRate:   *sampleRate,
	}, nil
}
