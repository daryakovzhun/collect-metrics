package filestore

import (
	"encoding/json"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
	"github.com/daryakovzhun/collect-metrics/internal/repository"
	"os"
	"sync"
)

type Config struct {
	Path string
}

type storage struct {
	cfg *Config
	mx  sync.RWMutex
}

func New(cfg *Config) repository.IFile {
	return &storage{cfg: cfg}
}

func (f *storage) Write(metrics []models.Metrics) error {
	if len(metrics) == 0 {
		return nil
	}

	existing, err := f.Read()
	if err != nil {
		return err
	}

	existingMap := make(map[string]models.Metrics, len(existing))
	for _, m := range existing {
		existingMap[m.ID] = m
	}

	for _, m := range metrics {
		existingMap[m.ID] = m
	}

	updated := make([]models.Metrics, 0, len(existingMap))
	for _, m := range existingMap {
		updated = append(updated, m)
	}

	data, err := json.Marshal(updated)
	if err != nil {
		return err
	}

	f.mx.Lock()
	defer f.mx.Unlock()
	return os.WriteFile(f.cfg.Path, data, 0666)
}

func (f *storage) Read() ([]models.Metrics, error) {
	f.mx.Lock()
	data, err := os.ReadFile(f.cfg.Path)
	f.mx.Unlock()
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
