package darwin

import (
	"fmt"
	"net"

	"github.com/abhinandpn/flux/internal/network"
)

// ListInterfaces discovers network interfaces using Go's standard library.
func ListInterfaces() ([]network.Interface, error) {
	systemInterfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("list macOS network interfaces: %w", err)
	}

	result := make([]network.Interface, 0, len(systemInterfaces))

	for _, iface := range systemInterfaces {
		item := network.Interface{
			Name:     iface.Name,
			Index:    iface.Index,
			Up:       iface.Flags&net.FlagUp != 0,
			Loopback: iface.Flags&net.FlagLoopback != 0,
		}

		if err := item.Validate(); err != nil {
			continue
		}

		result = append(result, item)
	}

	return result, nil
}
