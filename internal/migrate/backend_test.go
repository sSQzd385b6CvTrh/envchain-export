package migrate

import (
	"strings"
	"testing"
)

func TestBackendType_String(t *testing.T) {
	if BackendAkeyless.String() != "akeyless" {
		t.Errorf("expected 'akeyless', got %q", BackendAkeyless.String())
	}
}

func TestKnownBackends_ContainsAll(t *testing.T) {
	must := []BackendType{
		Backend1Password, BackendDoppler, BackendVault, BackendBitwarden,
		BackendAWSSecretsManager, BackendGCPSecretManager, BackendAzureKeyVault,
		BackendLastPass, BackendKeychain, BackendGopass, BackendInfisical,
		BackendHashiCorpVault, BackendSecretHub, BackendSSMParameterStore,
		BackendAkeyless, BackendDryRun,
	}
	set := make(map[BackendType]bool, len(KnownBackends))
	for _, b := range KnownBackends {
		set[b] = true
	}
	for _, b := range must {
		if !set[b] {
			t.Errorf("KnownBackends missing %q", b)
		}
	}
}

func TestParseBackend_Valid(t *testing.T) {
	b, err := ParseBackend("akeyless")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b != BackendAkeyless {
		t.Errorf("expected BackendAkeyless, got %v", b)
	}
}

func TestParseBackend_Invalid(t *testing.T) {
	_, err := ParseBackend("nonexistent-backend")
	if err == nil {
		t.Fatal("expected error for unknown backend")
	}
	if !strings.Contains(err.Error(), "unknown backend") {
		t.Errorf("unexpected error message: %v", err)
	}
}
