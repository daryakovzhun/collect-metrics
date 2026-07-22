package pg

import (
	"context"
	"fmt"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
	"github.com/daryakovzhun/collect-metrics/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	DatabaseDNS string
}

type database struct {
	pg *pgxpool.Pool
}

func New(ctx context.Context, cfg *Config) (repository.IRepository, error) {
	pool, err := pgxpool.New(ctx, cfg.DatabaseDNS)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	return &database{pg: pool}, nil
}

func (db *database) Close() {
	db.pg.Close()
}

func (db *database) Ping(ctx context.Context) error {
	err := db.pg.Ping(ctx)
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

func (db *database) SetGaugeMetric(metric models.Metrics) {
	//TODO implement me
	panic("implement me")
}

func (db *database) SetCounterMetric(metric models.Metrics) {
	//TODO implement me
	panic("implement me")
}

func (db *database) GetGaugeMetrics() ([]models.Metrics, error) {
	//TODO implement me
	panic("implement me")
}

func (db *database) GetCounterMetrics() ([]models.Metrics, error) {
	//TODO implement me
	panic("implement me")
}

func (db *database) GetAllMetrics() ([]models.Metrics, error) {
	//TODO implement me
	panic("implement me")
}

func (db *database) GetMetricByID(metric models.Metrics) (models.Metrics, error) {
	//TODO implement me
	panic("implement me")
}
