package migrate

import (
	"fmt"
	"strings"
)

// BackendType represents a supported secret backend.
type BackendType int

const (
	Backend1Password BackendType = iota
	BackendDoppler
	BackendVault
	BackendBitwarden
	BackendDryRun
)

// KnownBackends lists all supported backend identifiers.
var KnownBackends = []string{"1password", "doppler", "vault", "bitwarden", "dry-run"}

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
	case BackendDryRun:
		return "dry-run"
	default:
		return "unknown"
	}
}

// ParseBackend parses a backend string into a BackendType.
func ParseBackend(s string) (BackendType, error) {
	switch strings.ToLower(s) {
	case "1password":
		return Backend1Password, nil
	case "doppler":
		return BackendDoppler, nil
	case "vault":
		return BackendVault, nil
	case "bitwarden":
		return BackendBitwarden, nil
	case "dry-run":
		return BackendDryRun, nil
	default:
		return 0, fmt.Errorf("unknown backend %q: must be one of %s", s, strings.Join(KnownBackends, ", "))
	}
}
