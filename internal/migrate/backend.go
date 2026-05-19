package migrate

import (
	"fmt"
	"strings"
)

// BackendType identifies a supported secret backend.
type BackendType int

const (
	Backend1Password BackendType = iota
	BackendDoppler
	BackendVault
	BackendBitwarden
	BackendAWSSecretsManager
	BackendGCPSecretManager
)

// KnownBackends is the list of all supported backends.
var KnownBackends = []BackendType{
	Backend1Password,
	BackendDoppler,
	BackendVault,
	BackendBitwarden,
	BackendAWSSecretsManager,
	BackendGCPSecretManager,
}

func (b BackendType) String() string {
	switch b {
	case Backend1Password:
		return "1password"
	case BackendDoppler:
		return "doppler"
	case BackendVault:
		return "vault"
	case BackendBitwarden:
		return "bitwarden"
	case BackendAWSSecretsManager:
		return "awssecretsmanager"
	case BackendGCPSecretManager:
		return "gcpsecretmanager"
	default:
		return "unknown"
	}
}

// ParseBackend converts a string to a BackendType.
func ParseBackend(s string) (BackendType, error) {
	for _, b := range KnownBackends {
		if strings.EqualFold(b.String(), s) {
			return b, nil
		}
	}
	return 0, fmt.Errorf("unknown backend %q: must be one of %s", s, joinBackendNames())
}

func joinBackendNames() string {
	names := make([]string, len(KnownBackends))
	for i, b := range KnownBackends {
		names[i] = b.String()
	}
	return strings.Join(names, ", ")
}
