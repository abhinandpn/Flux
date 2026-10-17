package core

import (
	"errors"
	"math"
	"testing"
	"time"
)

func TestCalculate(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	previous := CounterSample{
		Interface:     "en0",
		BytesReceived: 1000,
		BytesSent:     2000,
		Timestamp:     start,
	}

	current := CounterSample{
		Interface:     "en0",
		BytesReceived: 501000,
		BytesSent:     102000,
		Timestamp:     start.Add(2 * time.Second),
	}

	got, err := Calculate(previous, current)
	if err != nil {
		t.Fatalf("Calculate() error = %v", err)
	}

	if got.Interface != "en0" {
		t.Errorf("Interface = %q, want en0", got.Interface)
	}

	if math.Abs(got.DownloadBPS-250000) > 0.001 {
		t.Errorf("DownloadBPS = %v, want 250000", got.DownloadBPS)
	}

	if math.Abs(got.UploadBPS-50000) > 0.001 {
		t.Errorf("UploadBPS = %v, want 50000", got.UploadBPS)
	}
}

func TestCalculateZeroTraffic(t *testing.T) {
	start := time.Now()

	previous := CounterSample{
		Interface: "en0",
		Timestamp: start,
	}

	current := CounterSample{
		Interface: "en0",
		Timestamp: start.Add(time.Second),
	}

	got, err := Calculate(previous, current)
	if err != nil {
		t.Fatalf("Calculate() error = %v", err)
	}

	if got.DownloadBPS != 0 || got.UploadBPS != 0 {
		t.Fatalf("expected zero rates, got %+v", got)
	}
}

func TestCalculateCounterReset(t *testing.T) {
	start := time.Now()

	previous := CounterSample{
		Interface:     "en0",
		BytesReceived: 5000,
		BytesSent:     3000,
		Timestamp:     start,
	}

	current := CounterSample{
		Interface:     "en0",
		BytesReceived: 100,
		BytesSent:     4000,
		Timestamp:     start.Add(time.Second),
	}

	_, err := Calculate(previous, current)
	if !errors.Is(err, ErrCounterReset) {
		t.Fatalf("error = %v, want ErrCounterReset", err)
	}
}

func TestCalculateInvalidInterval(t *testing.T) {
	now := time.Now()

	sample := CounterSample{
		Interface: "en0",
		Timestamp: now,
	}

	_, err := Calculate(sample, sample)
	if !errors.Is(err, ErrInvalidSample) {
		t.Fatalf("error = %v, want ErrInvalidSample", err)
	}
}

func TestCalculateInterfaceChanged(t *testing.T) {
	now := time.Now()

	previous := CounterSample{
		Interface: "en0",
		Timestamp: now,
	}

	current := CounterSample{
		Interface: "en1",
		Timestamp: now.Add(time.Second),
	}

	_, err := Calculate(previous, current)
	if !errors.Is(err, ErrInterfaceChanged) {
		t.Fatalf("error = %v, want ErrInterfaceChanged", err)
	}
}

func TestCalculateMissingTimestamp(t *testing.T) {
	previous := CounterSample{
		Interface: "en0",
		Timestamp: time.Now(),
	}

	current := CounterSample{
		Interface: "en0",
		Timestamp: time.Time{},
	}

	_, err := Calculate(previous, current)
	if !errors.Is(err, ErrInvalidSample) {
		t.Fatalf("error = %v, want ErrInvalidSample", err)
	}
}

func TestFormatRate(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  string
	}{
		{
			name:  "bytes",
			input: 500,
			want:  "0.00 Mbps",
		},
		{
			name:  "kibibytes",
			input: 2048,
			want:  "0.02 Mbps",
		},
		{
			name:  "mebibytes",
			input: 1024 * 1024,
			want:  "8.39 Mbps",
		},
		{
			name:  "zero",
			input: 0,
			want:  "0.00 Mbps",
		},
		{
			name:  "negative",
			input: -1,
			want:  "invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatRate(tt.input)
			if got != tt.want {
				t.Errorf("FormatRate() = %q, want %q", got, tt.want)
			}
		})
	}
}
