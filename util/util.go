package util

import (
	"context"
	"github.com/juju/errors"
	"github.com/ngaut/log"
	"time"
)

func RetryOnError(ctx context.Context, retryCount int, fn func() error) error {
	var err error
	for i := 0; i < retryCount; i++ {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		err = fn()
		if err == nil {
			break
		}

		log.Errorf("retry error %d times, %v", i+1, err)
		Sleep(ctx, 2*time.Second)
	}

	return errors.Trace(err)
}

// Sleep defines special `sleep` with context
func Sleep(ctx context.Context, sleepTime time.Duration) {
	ticker := time.NewTicker(sleepTime)
	defer ticker.Stop()

	select {
	case <-ctx.Done():
		return
	case <-ticker.C:
		return
	}
}
