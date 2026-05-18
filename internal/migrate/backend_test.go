package migrate

import (
	"testing"
)

func TestBackendType_String(t *testing.T) {
	cases := []struct {
		backend  BackendType
		expected string
	}{
		{Backend1Password, "1password"},
		{BackendDoppler, "doppler"},
		{BackendVault, "vault"},
	}
	for _, tc := range cases {
		if tc.backend.String() != tc.expected {
			t.Errorf("expected %q, got %q", tc.expected, tc.backend.String())
		}
	}
}

func TestKnownBackends_ContainsAll(t *testing.T) {
	expected := map[BackendType]bool{
		Backend1Password: false,
		BackendDoppler:   false,
		BackendVault:     false,
	}
	for _, b := range KnownBackends {
		expected[b] = true
	}
	for b, found := range expected {
		if !found {
			t.Errorf("KnownBackends missing %q", b)
		}
	}
}

func TestParseBackend_Valid(t *testing.T) {
	cases := []string{"1password", "doppler", "vault"}
	for _, s := range cases {
		b, err := ParseBackend(s)
		if err != nil {
			t.Errorf("unexpected error for %q: %v", s, err)
		}
		if b.String() != s {
			t.Errorf("expected %q, got %q", s, b.String())
		}
	}
}

func TestParseBackend_Invalid(t *testing.T) {
	_, err := ParseBackend("unknown-backend")
	if err == nil {
		t.Error("expected error for unknown backend, got nil")
	}
}
