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
		{BackendBitwarden, "bitwarden"},
		{BackendDryRun, "dry-run"},
	}

	for _, tc := range cases {
		if got := tc.backend.String(); got != tc.expected {
			t.Errorf("String() = %q, want %q", got, tc.expected)
		}
	}
}

func TestKnownBackends_ContainsAll(t *testing.T) {
	expected := []string{"1password", "doppler", "vault", "bitwarden", "dry-run"}
	if len(KnownBackends) != len(expected) {
		t.Fatalf("KnownBackends length = %d, want %d", len(KnownBackends), len(expected))
	}
	for i, name := range expected {
		if KnownBackends[i] != name {
			t.Errorf("KnownBackends[%d] = %q, want %q", i, KnownBackends[i], name)
		}
	}
}

func TestParseBackend_Valid(t *testing.T) {
	cases := []struct {
		input    string
		expected BackendType
	}{
		{"1password", Backend1Password},
		{"doppler", BackendDoppler},
		{"vault", BackendVault},
		{"bitwarden", BackendBitwarden},
		{"dry-run", BackendDryRun},
		{"DOPPLER", BackendDoppler},
		{"Vault", BackendVault},
	}

	for _, tc := range cases {
		got, err := ParseBackend(tc.input)
		if err != nil {
			t.Errorf("ParseBackend(%q) unexpected error: %v", tc.input, err)
			continue
		}
		if got != tc.expected {
			t.Errorf("ParseBackend(%q) = %v, want %v", tc.input, got, tc.expected)
		}
	}
}

func TestParseBackend_Invalid(t *testing.T) {
	_, err := ParseBackend("notabackend")
	if err == nil {
		t.Fatal("expected error for unknown backend, got nil")
	}
}
