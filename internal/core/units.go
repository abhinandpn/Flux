package core

import (
	"fmt"
	"math"
)

func FormatRate(bytesPerSecond float64) string {
	if math.IsNaN(bytesPerSecond) ||
		math.IsInf(bytesPerSecond, 0) ||
		bytesPerSecond < 0 {
		return "invalid"
	}

	megabitsPerSecond := bytesPerSecond * 8 / 1_000_000
	return fmt.Sprintf("%.2f Mbps", megabitsPerSecond)
}
