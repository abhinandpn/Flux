//go:build darwin

package darwin

import (
	"context"

	"github.com/abhinandpn/flux/internal/core"
)

// Provider adapts native macOS counters to the monitor interface.
type Provider struct{}

func (Provider) ReadCounters(
	ctx context.Context,
) ([]core.CounterSample, error) {
	return ReadCounters(ctx)
}
