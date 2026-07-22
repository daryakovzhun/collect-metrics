package app

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/daryakovzhun/collect-metrics/internal/utils"
)

type AgentConfig struct {
	ServerAddress  string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

type ServerConfig struct {
	Address         string `env:"ADDRESS"`
	StoreInterval   *int   `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         *bool  `env:"RESTORE"`
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

	if cfg.StoreInterval == nil {
		cfg.StoreInterval = utils.ToPointer(0)
		flag.IntVar(cfg.StoreInterval, "i", 300, "store interval")
	}

	if cfg.FileStoragePath == "" {
		flag.StringVar(&cfg.FileStoragePath, "f", "metrics.txt", "path to store files")
	}

	if cfg.Restore == nil {
		cfg.Restore = utils.ToPointer(false)
		flag.BoolVar(cfg.Restore, "r", false, "restore files")
	}

	flag.Parse()

	return &cfg, nil
}
