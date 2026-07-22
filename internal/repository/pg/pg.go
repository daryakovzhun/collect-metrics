package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
	"github.com/daryakovzhun/collect-metrics/internal/repository"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	db, err := sql.Open("postgres", cfg.DatabaseDNS)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to open sql.DB for migrations: %w", err)
	}
	defer db.Close()

	migrationsDir := "file://migrations"
	if err = startMigrations(db, migrationsDir); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &database{pg: pool}, nil
}

func startMigrations(db *sql.DB, migrationsDir string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationsDir,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed apply migrations: %w", err)
	}

	return nil
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

func (db *database) SetGaugeMetric(ctx context.Context, metric models.Metrics) error {
	return db.updateMetric(ctx, &metric)
}

func (db *database) SetCounterMetric(ctx context.Context, metric models.Metrics) error {
	return db.updateMetric(ctx, &metric)
}

func (db *database) GetAllMetrics(ctx context.Context) ([]models.Metrics, error) {
	rows, err := db.pg.Query(ctx, `SELECT name, type, delta, value, hash FROM metrics`)
	if err != nil {
		return nil, fmt.Errorf("failed to query get all metrics: %w", err)
	}
	defer rows.Close()

	var metrics []models.Metrics
	for rows.Next() {
		m := models.Metrics{}
		err = rows.Scan(&m.ID, &m.MType, &m.Delta, &m.Value, &m.Hash)
		if err != nil {
			return nil, fmt.Errorf("failed to scan metric row: %w", err)
		}

		metrics = append(metrics, m)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan metric rows: %w", err)
	}

	return metrics, nil
}

func (db *database) GetMetricByID(ctx context.Context, metric models.Metrics) (models.Metrics, error) {
	m := models.Metrics{}
	err := db.pg.QueryRow(ctx, `SELECT name, type, delta, value, hash FROM metrics where name = $1`, metric.ID).
		Scan(&m.ID, &m.MType, &m.Delta, &m.Value, &m.Hash)
	if err != nil {
		return models.Metrics{}, fmt.Errorf("failed to get metric by ID: %w", err)
	}

	return m, nil
}

func (db *database) updateMetric(ctx context.Context, metric *models.Metrics) error {
	_, err := db.pg.Exec(ctx, `
        INSERT INTO metrics (name, type, delta, value) 
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (name) 
        DO UPDATE SET
            type  = EXCLUDED.type,
            delta = EXCLUDED.delta,
            value = EXCLUDED.value,
            updated_at = NOW()`,
		metric.ID, metric.MType, metric.Delta, metric.Value)

	if err != nil {
		return fmt.Errorf("failed to update %s metric: %w", metric.MType, err)
	}

	return nil
}
