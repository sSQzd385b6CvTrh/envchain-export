package migrate

import "fmt"

// BackendType identifies a supported secret backend.
type BackendType string

const (
	Backend1Password BackendType = "1password"
	BackendDoppler   BackendType = "doppler"
	BackendVault     BackendType = "vault"
)

// String returns the string representation of the backend type.
func (b BackendType) String() string {
	return string(b)
}

// KnownBackends lists all supported backend identifiers.
var KnownBackends = []BackendType{
	Backend1Password,
	BackendDoppler,
	BackendVault,
}

// ParseBackend parses a string into a BackendType, returning an error if
// the value is not recognised.
func ParseBackend(s string) (BackendType, error) {
	for _, b := range KnownBackends {
		if string(b) == s {
			return b, nil
		}
	}
	return "", fmt.Errorf("unknown backend %q: must be one of 1password, doppler, vault", s)
}
