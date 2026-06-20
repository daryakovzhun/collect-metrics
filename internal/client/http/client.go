package httpclient

import (
	"fmt"
	"github.com/daryakovzhun/collect-metrics/internal/client"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
	"net/http"
	"time"
)

var (
	updateCounterEndpoint = "/update/%s/%s/%d"
	updateGaugeEndpoint   = "/update/%s/%s/%f"
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

func (h *httpClient) sendRequest(metric *models.Metrics) (*http.Response, error) {
	resURL, err := formURL(h.cfg.URL, metric)
	if err != nil {
		return nil, fmt.Errorf("failed to form url, err: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, resURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "text/plain")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	return resp, nil
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
