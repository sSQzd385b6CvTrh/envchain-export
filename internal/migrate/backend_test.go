package migrate

import (
	"strings"
	"testing"
)

func TestBackendType_String(t *testing.T) {
	if BackendPassbolt.String() != "passbolt" {
		t.Errorf("expected 'passbolt', got %q", BackendPassbolt.String())
	}
}

func TestKnownBackends_ContainsAll(t *testing.T) {
	expected := []BackendType{
		Backend1Password, BackendDoppler, BackendVault, BackendBitwarden,
		BackendAWSSecrets, BackendGCPSecret, BackendAzureKeyVault, BackendLastPass,
		BackendKeychain, BackendGopass, BackendInfisical, BackendHashiCorpVault,
		BackendSecretHub, BackendSSMParameter, BackendAkeyless, BackendKeePass,
		BackendAge, BackendPassbolt, BackendDryRun,
	}
	for _, e := range expected {
		found := false
		for _, k := range knownBackends {
			if k == e {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("backend %q not found in knownBackends", e)
		}
	}
}

func TestParseBackend_Valid(t *testing.T) {
	cases := []struct {
		input    string
		want     BackendType
	}{
		{"passbolt", BackendPassbolt},
		{"Passbolt", BackendPassbolt},
		{"PASSBOLT", BackendPassbolt},
		{"dryrun", BackendDryRun},
		{"1password", Backend1Password},
	}
	for _, tc := range cases {
		got, err := ParseBackend(tc.input)
		if err != nil {
			t.Errorf("ParseBackend(%q) unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("ParseBackend(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestParseBackend_Invalid(t *testing.T) {
	_, err := ParseBackend("nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown backend, got nil")
	}
	if !strings.Contains(err.Error(), "unknown backend") {
		t.Errorf("unexpected error message: %v", err)
	}
}
