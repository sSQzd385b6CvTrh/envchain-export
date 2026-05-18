package migrate

import (
	"testing"
)

func TestBackendType_String(t *testing.T) {
	tests := []struct {
		backend  BackendType
		expected string
	}{
		{Backend1Password, "1password"},
		{BackendDoppler, "doppler"},
	}

	for _, tt := range tests {
		if got := tt.backend.String(); got != tt.expected {
			t.Errorf("BackendType.String() = %q, want %q", got, tt.expected)
		}
	}
}

func TestKnownBackends_ContainsAll(t *testing.T) {
	known := KnownBackends()
	if len(known) < 2 {
		t.Fatalf("expected at least 2 known backends, got %d", len(known))
	}

	has := func(b BackendType) bool {
		for _, k := range known {
			if k == b {
				return true
			}
		}
		return false
	}

	for _, b := range []BackendType{Backend1Password, BackendDoppler} {
		if !has(b) {
			t.Errorf("KnownBackends() missing %q", b)
		}
	}
}

// Compile-time checks: ensure mock types satisfy the interfaces.
var _ SecretWriter = (*mockWriter)(nil)
var _ SecretReader = (*mockReader)(nil)

type mockWriter struct{}

func (m *mockWriter) WriteSecrets(_ string, _ map[string]string) error { return nil }

type mockReader struct{}

func (m *mockReader) ListNamespaces() ([]string, error)            { return nil, nil }
func (m *mockReader) ReadSecrets(_ string) (map[string]string, error) { return nil, nil }
