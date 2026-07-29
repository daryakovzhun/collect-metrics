package retryrepository

import (
	"context"
	"errors"
	"github.com/daryakovzhun/collect-metrics/internal/logger"
	models "github.com/daryakovzhun/collect-metrics/internal/model"
	"go.uber.org/zap"
	"net"
	"time"
)

var (
	delays = []time.Duration{1, 3, 5}
)

type retryableFunc func() error

func retry(ctx context.Context, fn retryableFunc) error {
	var lastErr error
	for attempt := 0; attempt < len(delays); attempt++ {
		// Проверяем отмену контекста перед каждой попыткой
		if ctx.Err() != nil {
			return ctx.Err()
		}

		err := fn()
		if err == nil {
			return nil
		}
		lastErr = err

		if errors.Is(err, models.ErrConnection) {
			logger.Log.Error("failed do func",
				zap.Int("attempt", attempt+1), zap.Error(lastErr))
			time.Sleep(delays[attempt] * time.Second)
			continue
		}

		return lastErr
	}
	return lastErr
}

// isNetworkError определяет ошибки соединения (из предыдущего ответа).
func isNetworkError(err error) bool {
	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}
	if errors.Is(err, net.ErrClosed) {
		return true
	}
	// При необходимости добавьте проверки специфических ошибок вашего драйвера БД,
	// например, для pgx: errors.Is(err, pgx.ErrDeadlineExceeded) или по тексту.
	return false
}
