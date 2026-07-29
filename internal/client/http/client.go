package httpclient

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/daryakovzhun/collect-metrics/internal/client"
	"github.com/daryakovzhun/collect-metrics/internal/logger"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/url"
	"time"
)

var (
	updateEndpoint  = "/update"
	updatesEndpoint = "/updates/"

	delays       = []time.Duration{1, 3, 5}
	countRetries = 3
)

type Config struct {
	Timeout time.Duration
	URL     string
}

type httpClient struct {
	cfg    *Config
	client *http.Client
}

func New(cfg *Config) client.IClient {
	cl := httpClient{
		cfg: cfg,
		client: &http.Client{
			Timeout:   cfg.Timeout,
			Transport: http.DefaultTransport.(*http.Transport).Clone(),
		},
	}

	return &cl
}

func (h *httpClient) SendMetric(ctx context.Context, metric *models.Metrics) error {
	resp, err := h.sendRequestWithRetries(ctx, updateEndpoint, metric)
	if err != nil {
		return fmt.Errorf("request error, err: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error status code: %d", resp.StatusCode)
	}

	return nil
}

func (h *httpClient) SendMetrics(ctx context.Context, metrics []models.Metrics) error {
	resp, err := h.sendRequestWithRetries(ctx, updatesEndpoint, metrics)
	if err != nil {
		return fmt.Errorf("request error, err: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error status code: %d", resp.StatusCode)
	}

	return nil
}

func (h *httpClient) sendRequestWithRetries(ctx context.Context, endpoint string, body interface{}) (*http.Response, error) {
	var lastErr error

	for i := 0; i < countRetries; i++ {
		resp, err := h.sendRequest(ctx, endpoint, body)
		if err == nil {
			return resp, nil
		}

		lastErr = err
		if isConnectionError(lastErr) {
			logger.Log.Error("failed send request",
				zap.Int("attempt", i+1), zap.Error(lastErr),
				zap.String("endpoint", endpoint))
			time.Sleep(delays[i] * time.Second)
			continue
		}

		break
	}

	return nil, lastErr
}

func (h *httpClient) sendRequest(ctx context.Context, endpoint string, body interface{}) (*http.Response, error) {
	resURL := h.cfg.URL + endpoint

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal metrics error, err: %w", err)
	}

	compressBody, err := compress(bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("compress metrics error, err: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, resURL, compressBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}

	return resp, nil
}

func compress(data []byte) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	gzWriter, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip writer: %w", err)
	}

	if _, err = gzWriter.Write(data); err != nil {
		return nil, fmt.Errorf("gzip write error: %w", err)
	}
	if err = gzWriter.Close(); err != nil {
		return nil, fmt.Errorf("gzip close error: %w", err)
	}

	return &buf, nil
}

func isConnectionError(err error) bool {
	if err == nil {
		return false
	}

	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		err = urlErr.Err
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}

	if errors.Is(err, net.ErrClosed) {
		return true
	}

	return false
}
