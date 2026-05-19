package migrate

import (
	"testing"
)

func TestBackendType_String(t *testing.T) {
	cases := []struct {
		backend BackendType
		want    string
	}{
		{Backend1Password, "1password"},
		{BackendDoppler, "doppler"},
		{BackendVault, "vault"},
		{BackendBitwarden, "bitwarden"},
		{BackendAWSSecretsManager, "awssecretsmanager"},
		{BackendGCPSecretManager, "gcpsecretmanager"},
	}
	for _, tc := range cases {
		if got := tc.backend.String(); got != tc.want {
			t.Errorf("BackendType(%d).String() = %q, want %q", tc.backend, got, tc.want)
		}
	}
}

func TestKnownBackends_ContainsAll(t *testing.T) {
	expected := []string{
		"1password", "doppler", "vault", "bitwarden", "awssecretsmanager", "gcpsecretmanager",
	}
	if len(KnownBackends) != len(expected) {
		t.Fatalf("KnownBackends length = %d, want %d", len(KnownBackends), len(expected))
	}
	for i, name := range expected {
		if KnownBackends[i].String() != name {
			t.Errorf("KnownBackends[%d] = %q, want %q", i, KnownBackends[i].String(), name)
		}
	}
}

func TestParseBackend_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  BackendType
	}{
		{"1password", Backend1Password},
		{"Doppler", BackendDoppler},
		{"VAULT", BackendVault},
		{"bitwarden", BackendBitwarden},
		{"awssecretsmanager", BackendAWSSecretsManager},
		{"gcpsecretmanager", BackendGCPSecretManager},
	}
	for _, tc := range cases {
		got, err := ParseBackend(tc.input)
		if err != nil {
			t.Errorf("ParseBackend(%q) unexpected error: %v", tc.input, err)
			continue
		}
		if got != tc.want {
			t.Errorf("ParseBackend(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestParseBackend_Invalid(t *testing.T) {
	_, err := ParseBackend("notabackend")
	if err == nil {
		t.Fatal("expected error for unknown backend, got nil")
	}
}
