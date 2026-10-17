package core

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidSample    = errors.New("invalid network sample")
	ErrCounterReset     = errors.New("network counter reset")
	ErrInterfaceChanged = errors.New("network interface changed")
)

// Calculate computes traffic rates between two counter samples.
func Calculate(previous, current CounterSample) (TrafficRate, error) {
	if previous.Interface == "" || current.Interface == "" {
		return TrafficRate{}, fmt.Errorf(
			"%w: interface name is empty", ErrInvalidSample,
		)
	}

	if previous.Interface != current.Interface {
		return TrafficRate{}, ErrInterfaceChanged
	}

	if previous.Timestamp.IsZero() || current.Timestamp.IsZero() {
		return TrafficRate{}, fmt.Errorf(
			"%w: timestamp is missing", ErrInvalidSample,
		)
	}

	elapsed := current.Timestamp.Sub(previous.Timestamp)
	if elapsed <= 0 {
		return TrafficRate{}, fmt.Errorf(
			"%w: elapsed time must be positive", ErrInvalidSample,
		)
	}

	if current.BytesReceived < previous.BytesReceived ||
		current.BytesSent < previous.BytesSent {
		return TrafficRate{}, ErrCounterReset
	}

	seconds := elapsed.Seconds()

	return TrafficRate{
		Interface:   current.Interface,
		DownloadBPS: float64(current.BytesReceived-previous.BytesReceived) / seconds,
		UploadBPS:   float64(current.BytesSent-previous.BytesSent) / seconds,
	}, nil
}
