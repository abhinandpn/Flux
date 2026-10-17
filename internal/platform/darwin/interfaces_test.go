package darwin

import "testing"

func TestListInterfaces(t *testing.T) {
	interfaces, err := ListInterfaces()
	if err != nil {
		t.Fatalf("ListInterfaces() error = %v", err)
	}

	if len(interfaces) == 0 {
		t.Fatal("expected at least one network interface")
	}

	for _, iface := range interfaces {
		if err := iface.Validate(); err != nil {
			t.Errorf("invalid interface %+v: %v", iface, err)
		}
	}
}
