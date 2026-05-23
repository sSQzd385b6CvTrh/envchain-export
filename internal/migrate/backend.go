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
	"github.com/envchain-export/internal/vault"
)

// BackendType represents a supported secret backend.
type BackendType int

const (
	Backend1Password BackendType = iota
	BackendDoppler
	BackendVault
	BackendAWSSecretsManager
	BackendGCPSecretManager
	BackendAzureKeyVault
	BackendLastPass
	BackendKeychain
	BackendGopass
	BackendBitwarden
	BackendInfisical
	BackendHashiCorpVault
)

var backendNames = map[BackendType]string{
	Backend1Password:         "1password",
	BackendDoppler:           "doppler",
	BackendVault:             "vault",
	BackendAWSSecretsManager: "aws-secrets-manager",
	BackendGCPSecretManager:  "gcp-secret-manager",
	BackendAzureKeyVault:     "azure-key-vault",
	BackendLastPass:          "lastpass",
	BackendKeychain:          "keychain",
	BackendGopass:            "gopass",
	BackendBitwarden:         "bitwarden",
	BackendInfisical:         "infisical",
	BackendHashiCorpVault:    "hashicorp-vault",
}

// String returns the string name of the backend.
func (b BackendType) String() string {
	if name, ok := backendNames[b]; ok {
		return name
	}
	return fmt.Sprintf("unknown(%d)", int(b))
}

// ParseBackend parses a string into a BackendType.
func ParseBackend(s string) (BackendType, error) {
	for bt, name := range backendNames {
		if strings.EqualFold(name, s) {
			return bt, nil
		}
	}
	return 0, fmt.Errorf("unknown backend %q; available: %s", s, joinBackendNames())
}

func joinBackendNames() string {
	names := make([]string, 0, len(backendNames))
	for _, name := range backendNames {
		names = append(names, name)
	}
	return strings.Join(names, ", ")
}

// SecretWriter is the interface that all backend writers must implement.
type SecretWriter interface {
	WriteSecret(namespace, key, value string) error
}

// NewBackendWriter constructs the appropriate SecretWriter for the given backend.
func NewBackendWriter(bt BackendType, opts map[string]string) (SecretWriter, error) {
	switch bt {
	case Backend1Password:
		return onepassword.NewWriter(opts["vault"]), nil
	case BackendDoppler:
		return doppler.NewWriter(opts["project"], opts["config"]), nil
	case BackendVault:
		return vault.NewWriter(opts["prefix"]), nil
	case BackendAWSSecretsManager:
		return awssecretsmanager.NewWriter(opts["region"]), nil
	case BackendGCPSecretManager:
		return gcpsecretmanager.NewWriter(opts["project"]), nil
	case BackendAzureKeyVault:
		return azurekeyvault.NewWriter(opts["vault"]), nil
	case BackendLastPass:
		return lastpass.NewWriter(opts["folder"]), nil
	case BackendKeychain:
		return keychain.NewWriter(), nil
	case BackendGopass:
		return gopass.NewWriter(opts["prefix"]), nil
	case BackendBitwarden:
		return bitwarden.NewWriter(opts["org"], opts["collection"]), nil
	case BackendInfisical:
		return infisical.NewWriter(opts["project"], opts["env"]), nil
	case BackendHashiCorpVault:
		return hashicorpvault.NewWriter(opts["mount"]), nil
	default:
		return nil, fmt.Errorf("unsupported backend: %s", bt)
	}
}
