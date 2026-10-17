package monitor

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/abhinandpn/flux/internal/core"
)

var (
	ErrMissingProvider = errors.New("counter provider is missing")
	ErrInvalidInterval = errors.New("sampling interval must be positive")
)

// CounterProvider supplies cumulative network counter samples.
type CounterProvider interface {
	ReadCounters(ctx context.Context) ([]core.CounterSample, error)
}

// Sampler reads counters from the configured provider.
type Sampler struct {
	Provider CounterProvider
	Interval time.Duration
}

// Sample reads one batch of counter samples.
func (s Sampler) Sample(
	ctx context.Context,
) ([]core.CounterSample, error) {
	if s.Provider == nil {
		return nil, ErrMissingProvider
	}

	if s.Interval <= 0 {
		return nil, ErrInvalidInterval
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	samples, err := s.Provider.ReadCounters(ctx)
	if err != nil {
		return nil, fmt.Errorf("read network counters: %w", err)
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return samples, nil
}
