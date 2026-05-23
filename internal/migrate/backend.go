package migrate

import (
	"fmt"
	"strings"

	"github.com/envchain-export/internal/awssecretsmanager"
	"github.com/envchain-export/internal/azurekeyvault"
	"github.com/envchain-export/internal/bitwarden"
	"github.com/envchain-export/internal/doppler"
	"github.com/envchain-export/internal/gcpsecretmanager"
	"github.com/envchain-export/internal/gopass"
	"github.com/envchain-export/internal/hashicorpvault"
	"github.com/envchain-export/internal/infisical"
	"github.com/envchain-export/internal/keychain"
	"github.com/envchain-export/internal/lastpass"
	"github.com/envchain-export/internal/onepassword"
	"github.com/envchain-export/internal/secrethub"
	"github.com/envchain-export/internal/ssmparameterstore"
	"github.com/envchain-export/internal/vault"
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
	BackendAzureKeyVault
	BackendLastPass
	BackendKeychain
	BackendGopass
	BackendInfisical
	BackendHashiCorpVault
	BackendSecretHub
	BackendSSMParameterStore
)

var backendNames = map[BackendType]string{
	Backend1Password:         "1password",
	BackendDoppler:           "doppler",
	BackendVault:             "vault",
	BackendBitwarden:         "bitwarden",
	BackendAWSSecretsManager: "aws-secrets-manager",
	BackendGCPSecretManager:  "gcp-secret-manager",
	BackendAzureKeyVault:     "azure-key-vault",
	BackendLastPass:          "lastpass",
	BackendKeychain:          "keychain",
	BackendGopass:            "gopass",
	BackendInfisical:         "infisical",
	BackendHashiCorpVault:    "hashicorp-vault",
	BackendSecretHub:         "secrethub",
	BackendSSMParameterStore: "ssm-parameter-store",
}

func (b BackendType) String() string {
	if name, ok := backendNames[b]; ok {
		return name
	}
	return "unknown"
}

// KnownBackends returns all registered backend types.
func KnownBackends() []BackendType {
	backends := make([]BackendType, 0, len(backendNames))
	for b := range backendNames {
		backends = append(backends, b)
	}
	return backends
}

// ParseBackend parses a backend name string into a BackendType.
func ParseBackend(name string) (BackendType, error) {
	for b, n := range backendNames {
		if strings.EqualFold(n, name) {
			return b, nil
		}
	}
	return 0, fmt.Errorf("unknown backend %q; available: %s", name, joinBackendNames())
}

func joinBackendNames() string {
	names := make([]string, 0, len(backendNames))
	for _, n := range backendNames {
		names = append(names, n)
	}
	return strings.Join(names, ", ")
}

// SecretWriter is the interface all backend writers must satisfy.
type SecretWriter interface {
	WriteSecret(namespace, key, value string) error
}

// NewBackendWriter constructs a SecretWriter for the given backend using flags.
func NewBackendWriter(b BackendType, flags map[string]string) (SecretWriter, error) {
	switch b {
	case Backend1Password:
		return onepassword.NewWriter(flags["vault"]), nil
	case BackendDoppler:
		return doppler.NewWriter(flags["project"], flags["config"]), nil
	case BackendVault:
		return vault.NewWriter(flags["mount"]), nil
	case BackendBitwarden:
		return bitwarden.NewWriter(flags["org"], flags["collection"]), nil
	case BackendAWSSecretsManager:
		return awssecretsmanager.NewWriter(flags["prefix"]), nil
	case BackendGCPSecretManager:
		return gcpsecretmanager.NewWriter(flags["project"]), nil
	case BackendAzureKeyVault:
		return azurekeyvault.NewWriter(flags["vault"]), nil
	case BackendLastPass:
		return lastpass.NewWriter(flags["folder"]), nil
	case BackendKeychain:
		return keychain.NewWriter(), nil
	case BackendGopass:
		return gopass.NewWriter(flags["store"]), nil
	case BackendInfisical:
		return infisical.NewWriter(flags["project"], flags["env"]), nil
	case BackendHashiCorpVault:
		return hashicorpvault.NewWriter(flags["addr"], flags["mount"]), nil
	case BackendSecretHub:
		return secrethub.NewWriter(flags["org"], flags["repo"]), nil
	case BackendSSMParameterStore:
		return ssmparameterstore.NewWriter(flags["path"], flags["kms-key-id"]), nil
	default:
		return nil, fmt.Errorf("no writer registered for backend %s", b)
	}
}
