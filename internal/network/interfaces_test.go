package network

import (
	"errors"
	"testing"
)

func TestInterfaceValidate(t *testing.T) {
	tests := []struct {
		name    string
		iface   Interface
		wantErr error
	}{
		{
			name: "valid interface",
			iface: Interface{
				Name:  "en0",
				Index: 4,
				Up:    true,
			},
		},
		{
			name: "empty name",
			iface: Interface{
				Index: 4,
			},
			wantErr: ErrEmptyInterfaceName,
		},
		{
			name: "invalid index",
			iface: Interface{
				Name:  "en0",
				Index: 0,
			},
			wantErr: ErrInvalidInterfaceID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.iface.Validate()

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Validate() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
