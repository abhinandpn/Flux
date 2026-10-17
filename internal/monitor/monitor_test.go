package monitor

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/abhinandpn/flux/internal/core"
)

type fakeProvider struct {
	samples []core.CounterSample
	err     error
}

func (f *fakeProvider) ReadCounters(
	ctx context.Context,
) ([]core.CounterSample, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.samples, nil
}

func TestSamplerSample(t *testing.T) {
	provider := &fakeProvider{
		samples: []core.CounterSample{
			{
				Interface:     "en0",
				BytesReceived: 1000,
				BytesSent:     2000,
				Timestamp:     time.Now(),
			},
		},
	}

	sampler := Sampler{
		Provider: provider,
		Interval: time.Second,
	}

	samples, err := sampler.Sample(context.Background())
	if err != nil {
		t.Fatalf("Sample() error = %v", err)
	}

	if len(samples) != 1 {
		t.Fatalf("got %d samples, want 1", len(samples))
	}

	if samples[0].Interface != "en0" {
		t.Errorf("interface = %q, want en0", samples[0].Interface)
	}
}

func TestSamplerMissingProvider(t *testing.T) {
	sampler := Sampler{Interval: time.Second}

	_, err := sampler.Sample(context.Background())
	if !errors.Is(err, ErrMissingProvider) {
		t.Fatalf("error = %v, want ErrMissingProvider", err)
	}
}

func TestSamplerInvalidInterval(t *testing.T) {
	sampler := Sampler{
		Provider: &fakeProvider{},
		Interval: 0,
	}

	_, err := sampler.Sample(context.Background())
	if !errors.Is(err, ErrInvalidInterval) {
		t.Fatalf("error = %v, want ErrInvalidInterval", err)
	}
}

func TestSamplerProviderError(t *testing.T) {
	providerErr := errors.New("simulated read failure")

	sampler := Sampler{
		Provider: &fakeProvider{err: providerErr},
		Interval: time.Second,
	}

	_, err := sampler.Sample(context.Background())
	if !errors.Is(err, providerErr) {
		t.Fatalf("error = %v, want provider error", err)
	}
}

func TestMonitorCancellation(t *testing.T) {
	provider := &fakeProvider{
		samples: []core.CounterSample{
			{
				Interface: "en0",
				Timestamp: time.Now(),
			},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	updates := make(chan Update)

	done := make(chan error, 1)

	go func() {
		done <- (Monitor{
			Sampler: Sampler{
				Provider: provider,
				Interval: 10 * time.Millisecond,
			},
		}).Run(ctx, updates)
	}()

	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Run() error = %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("monitor did not stop after cancellation")
	}
}
