//go:build !darwin

package darwin

import (
	"context"
	"errors"

	"github.com/abhinandpn/flux/internal/core"
)

type Provider struct{}

func (Provider) ReadCounters(
	ctx context.Context,
) ([]core.CounterSample, error) {
	return nil, errors.New("macOS counter provider is unavailable on this OS")
}
