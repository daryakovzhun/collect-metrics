package filestore

import (
	"encoding/json"
	"fmt"
	"github.com/daryakovzhun/collect-metrics/internal/logger"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
	"github.com/daryakovzhun/collect-metrics/internal/repository"
	"go.uber.org/zap"
	"os"
	"path/filepath"
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
	err := ensureFileExists(cfg.Path)
	if err != nil {
		logger.Log.Error("failed to ensure filesystem directory", zap.String("path", cfg.Path), zap.Error(err))
	}

	return &storage{cfg: cfg}
}

// ensureFileExists проверяет существование файла, создаёт все необходимые папки
// и сам файл, если он отсутствует.
// Возвращает ошибку, если что-то пошло не так.
func ensureFileExists(path string) error {
	// 1. Проверяем, существует ли файл
	_, err := os.Stat(path)
	if err == nil {
		fmt.Println("Файл уже существует")
		return nil
	}

	if os.IsNotExist(err) {
		// 2. Файла нет – выделяем директорию из пути
		dir := filepath.Dir(path)

		// Если директория не пустая и не текущая, создаём её рекурсивно
		if dir != "" && dir != "." {
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				return fmt.Errorf("не удалось создать директории: %w", err)
			}
			fmt.Printf("Директории %s созданы\n", dir)
		}

		// 3. Создаём сам файл
		file, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("не удалось создать файл: %w", err)
		}
		defer file.Close() // закрываем сразу, файл пуст

		fmt.Printf("Файл %s создан\n", path)
		return nil
	}

	// Прочие ошибки (например, нет прав доступа к папке)
	return fmt.Errorf("ошибка при проверке файла: %w", err)
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
