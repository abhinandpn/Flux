//go:build !darwin

package darwin

import (
	"context"
	"errors"

	"github.com/abhinandpn/flux/internal/core"
)

var ErrCountersNotImplemented = errors.New(
	"macOS network counter collection is not implemented yet",
)

// ReadCounters will return cumulative byte counters for each
// supported network interface.
//
// The native macOS counter implementation will replace this stub.
func ReadCounters(
	ctx context.Context,
) ([]core.CounterSample, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return nil, ErrCountersNotImplemented
}
