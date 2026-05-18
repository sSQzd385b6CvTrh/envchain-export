package migrate

// SecretWriter is the interface that wraps secret writing operations.
// Any backend (1Password, Doppler, etc.) must implement this to be used
// as a migration target.
type SecretWriter interface {
	// WriteSecrets writes all key-value pairs for a given namespace/group.
	WriteSecrets(namespace string, secrets map[string]string) error
}

// SecretReader is the interface for reading secrets from a source backend.
type SecretReader interface {
	// ListNamespaces returns all available namespace names.
	ListNamespaces() ([]string, error)

	// ReadSecrets returns all key-value pairs for a given namespace.
	ReadSecrets(namespace string) (map[string]string, error)
}

// BackendType enumerates supported destination backends.
type BackendType string

const (
	Backend1Password BackendType = "1password"
	BackendDoppler   BackendType = "doppler"
)

// String returns the string representation of a BackendType.
func (b BackendType) String() string {
	return string(b)
}

// KnownBackends returns all supported backend identifiers.
func KnownBackends() []BackendType {
	return []BackendType{Backend1Password, BackendDoppler}
}
