package monitor

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/abhinandpn/flux/internal/core"
)

type Update struct {
	Rate core.TrafficRate
	Err  error
}

type Monitor struct {
	Sampler   Sampler
	Interface string
}

func (m Monitor) Run(ctx context.Context, updates chan<- Update) error {
	if updates == nil {
		return errors.New("updates channel is nil")
	}

	if m.Sampler.Provider == nil {
		return ErrMissingProvider
	}

	if m.Sampler.Interval <= 0 {
		return ErrInvalidInterval
	}

	previous, err := m.Sampler.Sample(ctx)
	if err == nil {
		previous, err = m.selectSamples(previous)
	}

	if err != nil {
		if ctx.Err() != nil {
			return nil
		}

		if !sendUpdate(ctx, updates, Update{Err: err}) {
			return nil
		}

		previous = nil
	}

	ticker := time.NewTicker(m.Sampler.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil

		case <-ticker.C:
			current, err := m.Sampler.Sample(ctx)
			if err == nil {
				current, err = m.selectSamples(current)
			}

			if err != nil {
				if ctx.Err() != nil {
					return nil
				}

				if !sendUpdate(ctx, updates, Update{Err: err}) {
					return nil
				}

				previous = nil
				continue
			}

			if previous != nil {
				for _, sample := range current {
					old, found := findSample(previous, sample.Interface)
					if !found {
						continue
					}

					rate, err := core.Calculate(old, sample)
					if err != nil {
						if !sendUpdate(ctx, updates, Update{
							Err: fmt.Errorf(
								"interface %s: %w",
								sample.Interface,
								err,
							),
						}) {
							return nil
						}
						continue
					}

					if !sendUpdate(ctx, updates, Update{Rate: rate}) {
						return nil
					}
				}
			}

			previous = current
		}
	}
}

func (m Monitor) selectSamples(
	samples []core.CounterSample,
) ([]core.CounterSample, error) {
	if m.Interface == "" {
		return samples, nil
	}

	for _, sample := range samples {
		if sample.Interface == m.Interface {
			return []core.CounterSample{sample}, nil
		}
	}

	return nil, fmt.Errorf(
		"interface %q not found in counter samples",
		m.Interface,
	)
}

func findSample(
	samples []core.CounterSample,
	name string,
) (core.CounterSample, bool) {
	for _, sample := range samples {
		if sample.Interface == name {
			return sample, true
		}
	}

	return core.CounterSample{}, false
}

func sendUpdate(
	ctx context.Context,
	updates chan<- Update,
	update Update,
) bool {
	select {
	case <-ctx.Done():
		return false
	case updates <- update:
		return true
	}
}
