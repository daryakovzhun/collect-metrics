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
	Key            string `env:"KEY"`
}

type ServerConfig struct {
	Address         string `env:"ADDRESS"`
	StoreInterval   *int   `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         *bool  `env:"RESTORE"`
	DB              string `env:"DATABASE_DSN"`
	Key             string `env:"KEY"`
}

func getAgentConfig() (*AgentConfig, error) {
	var cfg AgentConfig

	addr := flag.String("a", "localhost:8080", "address and port to run server")
	reportInterval := flag.Int("r", 10, "frequency of sending metrics to the server")
	pollInterval := flag.Int("p", 2, "the frequency of polling metrics from the package")
	key := flag.String("k", "", "the key to use for encryption")

	flag.Parse()

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	if cfg.ServerAddress == "" {
		cfg.ServerAddress = *addr
	}
	if cfg.ReportInterval == 0 {
		cfg.ReportInterval = *reportInterval
	}
	if cfg.PollInterval == 0 {
		cfg.PollInterval = *pollInterval
	}
	if cfg.Key == "" {
		cfg.Key = *key
	}

	return &cfg, nil
}

func getServerConfig() (*ServerConfig, error) {
	var cfg ServerConfig

	addr := flag.String("a", "localhost:8080", "address and port to run server")
	storeInterval := flag.Int("i", 300, "store interval")
	filePath := flag.String("f", "metrics.txt", "path to store files")
	restore := flag.Bool("r", false, "restore metrics")
	db := flag.String("d", "", "database connection string")
	key := flag.String("k", "", "the key to use for encryption")

	flag.Parse()

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	if cfg.Address == "" {
		cfg.Address = *addr
	}
	if cfg.StoreInterval == nil {
		cfg.StoreInterval = utils.ToPointer(*storeInterval)
	}
	if cfg.FileStoragePath == "" {
		cfg.FileStoragePath = *filePath
	}
	if cfg.Restore == nil {
		cfg.Restore = utils.ToPointer(*restore)
	}
	if cfg.DB == "" {
		cfg.DB = *db
	}
	if cfg.Key == "" {
		cfg.Key = *key
	}

	return &cfg, nil
}
