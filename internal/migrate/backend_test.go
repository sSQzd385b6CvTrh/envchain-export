package migrate

import (
	"testing"
)

func TestBackendType_String(t *testing.T) {
	tests := []struct {
		bt   BackendType
		want string
	}{
		{Backend1Password, "1password"},
		{BackendDoppler, "doppler"},
		{BackendEnpass, "enpass"},
		{BackendBitwarden, "bitwarden"},
	}
	for _, tt := range tests {
		if got := tt.bt.String(); got != tt.want {
			t.Errorf("BackendType(%d).String() = %q, want %q", int(tt.bt), got, tt.want)
		}
	}
}

func TestKnownBackends_ContainsAll(t *testing.T) {
	expected := []string{
		"1password", "doppler", "vault", "aws-secrets-manager",
		"gcp-secret-manager", "azure-key-vault", "lastpass", "keychain",
		"gopass", "infisical", "hashicorp-vault", "secrethub",
		"ssm-parameter-store", "akeyless", "keepass", "age",
		"passbolt", "conjur", "bitwarden", "enpass",
	}
	for _, name := range expected {
		if _, err := ParseBackend(name); err != nil {
			t.Errorf("ParseBackend(%q) returned error: %v", name, err)
		}
	}
}

func TestParseBackend_Valid(t *testing.T) {
	bt, err := ParseBackend("enpass")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bt != BackendEnpass {
		t.Errorf("got %v, want BackendEnpass", bt)
	}
}

func TestParseBackend_Invalid(t *testing.T) {
	_, err := ParseBackend("notabackend")
	if err == nil {
		t.Fatal("expected error for unknown backend, got nil")
	}
}
