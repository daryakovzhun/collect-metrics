package httpclient

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/daryakovzhun/collect-metrics/internal/client"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
	"net/http"
	"time"
)

var (
	updateCounterEndpoint = "/update/%s/%s/%d"
	updateGaugeEndpoint   = "/update/%s/%s/%f"
	updateEndpoint        = "/update"
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

func (h *httpClient) SendMetric(metric *models.Metrics) error {
	resp, err := h.sendRequest(metric)
	if err != nil {
		return fmt.Errorf("request error, err: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error status code: %d", resp.StatusCode)
	}

	return nil
}

func (h *httpClient) sendRequest(metrics *models.Metrics) (*http.Response, error) {
	resURL := h.cfg.URL + updateEndpoint

	body, err := json.Marshal(metrics)
	if err != nil {
		return nil, fmt.Errorf("marshal metrics error, err: %w", err)
	}

	compressBody, err := compress(body)
	if err != nil {
		return nil, fmt.Errorf("compress metrics error, err: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, resURL, compressBody)
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

func formURL(url string, metrics *models.Metrics) (string, error) {
	switch metrics.MType {
	case models.Counter:
		return url + fmt.Sprintf(updateCounterEndpoint, metrics.MType, metrics.ID, *metrics.Delta), nil
	case models.Gauge:
		return url + fmt.Sprintf(updateGaugeEndpoint, metrics.MType, metrics.ID, *metrics.Value), nil
	default:
		return "", fmt.Errorf("unsupported metrics type: %s", metrics.MType)
	}
}
