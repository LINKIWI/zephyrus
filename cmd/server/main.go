package main

import (
	"flag"
	"log"

	"zephyrus/internal/device"
	"zephyrus/internal/server"
)

type Config struct {
	Port       int
	Identifier string
}

func main() {
	cfg, err := parseConfig()
	if err != nil {
		panic(err)
	}

	log.Printf("main: using configuration: port=%d device=%s", cfg.Port, cfg.Identifier)

	log.Printf("main: finding and initializing device")
	temper, err := device.NewTemperClient(cfg.Identifier)
	if err != nil {
		panic(err)
	}

	if err := temper.Open(); err != nil {
		panic(err)
	}
	defer temper.Close()

	log.Printf("main: initializing Zephyrus gRPC server")
	zephyrus, err := server.NewZephyrusServer(temper)
	if err != nil {
		panic(err)
	}

	log.Printf("main: serving on port %d", cfg.Port)
	if err := zephyrus.Serve(cfg.Port); err != nil {
		panic(err)
	}
}

func parseConfig() (*Config, error) {
	port := flag.Int("port", 6840, "TCP port on which the gRPC server should listen")
	identifier := flag.String(
		"identifier",
		"temper",
		"Name used to uniquely identify the device associated with this server",
	)
	flag.Parse()

	return &Config{
		Port:       *port,
		Identifier: *identifier,
	}, nil
}
