package filestore

import (
	"encoding/json"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
	"github.com/daryakovzhun/collect-metrics/internal/repository"
	"os"
)

type Config struct {
	Path string
}

type storage struct {
	cfg *Config
}

func New(cfg *Config) repository.IFile {
	return &storage{cfg: cfg}
}

func (f *storage) Write(metrics []models.Metrics) error {
	if len(metrics) == 0 {
		return nil
	}

	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	return os.WriteFile(f.cfg.Path, data, 0666)
}

func (f *storage) Read() ([]models.Metrics, error) {
	data, err := os.ReadFile(f.cfg.Path)
	if err != nil {
		return nil, err
	}

	var metrics []models.Metrics
	if len(data) == 0 {
		return metrics, nil
	}

	if err = json.Unmarshal(data, &metrics); err != nil {
		return nil, err
	}

	return metrics, nil
}
