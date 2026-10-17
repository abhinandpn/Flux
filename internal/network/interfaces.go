package network

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyInterfaceName = errors.New("interface name is empty")
	ErrInvalidInterfaceID = errors.New("interface index must be positive")
)

// Interface describes a network interface detected by the OS.
type Interface struct {
	Name     string
	Index    int
	Up       bool
	Loopback bool
}

// Validate checks whether the interface has valid identity fields.
func (i Interface) Validate() error {
	if i.Name == "" {
		return ErrEmptyInterfaceName
	}

	if i.Index <= 0 {
		return fmt.Errorf(
			"%w: got %d", ErrInvalidInterfaceID, i.Index,
		)
	}

	return nil
}
