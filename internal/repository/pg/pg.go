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
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	DatabaseDNS string
}

type database struct {
	pool *pgxpool.Pool
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

	return &database{pool: pool}, nil
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
	db.pool.Close()
}

func (db *database) Ping(ctx context.Context) error {
	err := db.pool.Ping(ctx)
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
	rows, err := db.pool.Query(ctx, getAllMetrics)
	if err != nil {
		return nil, handleError("failed to query get all metrics", err)
	}
	defer rows.Close()

	var metrics []models.Metrics
	for rows.Next() {
		m := models.Metrics{}
		err = rows.Scan(&m.ID, &m.MType, &m.Delta, &m.Value, &m.Hash)
		if err != nil {
			return nil, handleError("failed to scan metric row", err)
		}

		metrics = append(metrics, m)
	}

	if err = rows.Err(); err != nil {
		return nil, handleError("failed to scan metric rows", err)
	}

	return metrics, nil
}

func (db *database) GetMetricByID(ctx context.Context, metric models.Metrics) (models.Metrics, error) {
	m := models.Metrics{}
	err := db.pool.QueryRow(ctx, getMetricByID, metric.ID).
		Scan(&m.ID, &m.MType, &m.Delta, &m.Value, &m.Hash)
	if err != nil {
		return models.Metrics{}, handleError("failed to get metric by ID", err)
	}

	return m, nil
}

func (db *database) updateMetric(ctx context.Context, metric *models.Metrics) error {
	_, err := db.pool.Exec(ctx, updateMetricsQuery,
		metric.ID, metric.MType, metric.Delta, metric.Value, metric.Hash)

	if err != nil {
		return handleError("failed to update metric", err)
	}

	return nil
}

func (db *database) UpdateMetrics(ctx context.Context, metrics []models.Metrics) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return handleError("failed to begin tx", err)
	}

	defer tx.Rollback(ctx)

	for _, metric := range metrics {
		_, err = tx.Exec(ctx, updateMetricsQuery, metric.ID, metric.MType,
			metric.Delta, metric.Value, metric.Hash)
		if err != nil {
			return handleError("failed to update metric", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return handleError("failed to commit tx", err)
	}

	return nil
}

func handleError(msg string, err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("msg: %s, err: %w", msg, models.ErrNotFound)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
		return fmt.Errorf("msg: %s, %w, err: %w", msg, models.ErrConnection, err)
	}

	return fmt.Errorf("msg: %s, err: %w", msg, err)
}
