package migrate

import (
	"fmt"
	"strings"

	"github.com/yourusername/envchain-export/internal/awssecretsmanager"
	"github.com/yourusername/envchain-export/internal/azurekeyvault"
	"github.com/yourusername/envchain-export/internal/bitwarden"
	"github.com/yourusername/envchain-export/internal/doppler"
	"github.com/yourusername/envchain-export/internal/gcpsecretmanager"
	"github.com/yourusername/envchain-export/internal/gopass"
	"github.com/yourusername/envchain-export/internal/infisical"
	"github.com/yourusername/envchain-export/internal/keychain"
	"github.com/yourusername/envchain-export/internal/lastpass"
	"github.com/yourusername/envchain-export/internal/onepassword"
	"github.com/yourusername/envchain-export/internal/vault"
)

// BackendType identifies a supported secret backend.
type BackendType string

const (
	Backend1Password      BackendType = "1password"
	BackendDoppler        BackendType = "doppler"
	BackendVault          BackendType = "vault"
	BackendBitwarden      BackendType = "bitwarden"
	BackendAWSSecrets     BackendType = "aws-secrets-manager"
	BackendGCPSecrets     BackendType = "gcp-secret-manager"
	BackendAzureKeyVault  BackendType = "azure-key-vault"
	BackendLastPass       BackendType = "lastpass"
	BackendKeychain       BackendType = "keychain"
	BackendGopass         BackendType = "gopass"
	BackendInfisical      BackendType = "infisical"
)

func (b BackendType) String() string {
	return string(b)
}

// KnownBackends lists all supported backend types.
var KnownBackends = []BackendType{
	Backend1Password,
	BackendDoppler,
	BackendVault,
	BackendBitwarden,
	BackendAWSSecrets,
	BackendGCPSecrets,
	BackendAzureKeyVault,
	BackendLastPass,
	BackendKeychain,
	BackendGopass,
	BackendInfisical,
}

// ParseBackend parses a string into a BackendType.
func ParseBackend(s string) (BackendType, error) {
	for _, b := range KnownBackends {
		if strings.EqualFold(s, string(b)) {
			return b, nil
		}
	}
	return "", fmt.Errorf("unknown backend %q; supported: %s", s, joinBackendNames())
}

func joinBackendNames() string {
	names := make([]string, len(KnownBackends))
	for i, b := range KnownBackends {
		names[i] = string(b)
	}
	return strings.Join(names, ", ")
}

// SecretWriter is the interface all backend writers must satisfy.
type SecretWriter interface {
	WriteSecret(namespace, key, value string) error
}

// NewBackendWriter constructs the appropriate writer for the given backend using opts.
func NewBackendWriter(b BackendType, opts map[string]string) (SecretWriter, error) {
	switch b {
	case Backend1Password:
		return onepassword.NewWriter(opts["vault"]), nil
	case BackendDoppler:
		return doppler.NewWriter(opts["project"], opts["config"]), nil
	case BackendVault:
		return vault.NewWriter(opts["mount"]), nil
	case BackendBitwarden:
		return bitwarden.NewWriter(opts["org"], opts["collection"]), nil
	case BackendAWSSecrets:
		return awssecretsmanager.NewWriter(opts["region"]), nil
	case BackendGCPSecrets:
		return gcpsecretmanager.NewWriter(opts["project"]), nil
	case BackendAzureKeyVault:
		return azurekeyvault.NewWriter(opts["vault"]), nil
	case BackendLastPass:
		return lastpass.NewWriter(opts["folder"]), nil
	case BackendKeychain:
		return keychain.NewWriter(), nil
	case BackendGopass:
		return gopass.NewWriter(opts["store"]), nil
	case BackendInfisical:
		return infisical.NewWriter(opts["project"], opts["environment"]), nil
	}
	return nil, fmt.Errorf("no writer implemented for backend %q", b)
}
