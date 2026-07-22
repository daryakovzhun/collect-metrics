package app

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

type AgentConfig struct {
	ServerAddress  string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

type ServerConfig struct {
	Address string `env:"ADDRESS"`
}

func getAgentConfig() (*AgentConfig, error) {
	var cfg AgentConfig

	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	if cfg.ServerAddress == "" {
		flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "address and port to run server")
	}

	if cfg.ReportInterval == 0 {
		flag.IntVar(&cfg.ReportInterval, "r", 10, "frequency of sending metrics to the server")
	}

	if cfg.PollInterval == 0 {
		flag.IntVar(&cfg.PollInterval, "p", 2, "the frequency of polling metrics from the package")
	}

	flag.Parse()

	return &cfg, nil
}

func getServerConfig() (*ServerConfig, error) {
	var cfg ServerConfig

	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	if cfg.Address == "" {
		flag.StringVar(&cfg.Address, "a", "localhost:8080", "address and port to run server")
	}

	flag.Parse()

	return &cfg, nil
}
