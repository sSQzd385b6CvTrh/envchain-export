package migrate

import (
	"fmt"
	"strings"

	"github.com/envchain-export/internal/akeyless"
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
type BackendType string

func (b BackendType) String() string { return string(b) }

const (
	Backend1Password        BackendType = "1password"
	BackendDoppler          BackendType = "doppler"
	BackendVault            BackendType = "vault"
	BackendBitwarden        BackendType = "bitwarden"
	BackendAWSSecretsManager BackendType = "aws-secrets-manager"
	BackendGCPSecretManager BackendType = "gcp-secret-manager"
	BackendAzureKeyVault    BackendType = "azure-key-vault"
	BackendLastPass         BackendType = "lastpass"
	BackendKeychain         BackendType = "keychain"
	BackendGopass           BackendType = "gopass"
	BackendInfisical        BackendType = "infisical"
	BackendHashiCorpVault   BackendType = "hashicorp-vault"
	BackendSecretHub        BackendType = "secrethub"
	BackendSSMParameterStore BackendType = "ssm-parameter-store"
	BackendAkeyless         BackendType = "akeyless"
	BackendDryRun           BackendType = "dry-run"
)

// KnownBackends lists all supported backend identifiers.
var KnownBackends = []BackendType{
	Backend1Password,
	BackendDoppler,
	BackendVault,
	BackendBitwarden,
	BackendAWSSecretsManager,
	BackendGCPSecretManager,
	BackendAzureKeyVault,
	BackendLastPass,
	BackendKeychain,
	BackendGopass,
	BackendInfisical,
	BackendHashiCorpVault,
	BackendSecretHub,
	BackendSSMParameterStore,
	BackendAkeyless,
	BackendDryRun,
}

// SecretWriter is the interface all backend writers must satisfy.
type SecretWriter interface {
	WriteSecret(namespace, key, value string) error
}

// ParseBackend parses a string into a known BackendType.
func ParseBackend(s string) (BackendType, error) {
	for _, b := range KnownBackends {
		if string(b) == s {
			return b, nil
		}
	}
	return "", fmt.Errorf("unknown backend %q; known backends: %s", s, joinBackendNames())
}

func joinBackendNames() string {
	names := make([]string, len(KnownBackends))
	for i, b := range KnownBackends {
		names[i] = string(b)
	}
	return strings.Join(names, ", ")
}

// NewBackendWriter constructs the SecretWriter for the given BackendType.
func NewBackendWriter(b BackendType, opts map[string]string) (SecretWriter, error) {
	switch b {
	case Backend1Password:
		return onepassword.NewWriter(opts["vault"]), nil
	case BackendDoppler:
		return doppler.NewWriter(opts["project"], opts["config"]), nil
	case BackendVault:
		return vault.NewWriter(opts["path"]), nil
	case BackendBitwarden:
		return bitwarden.NewWriter(opts["org"], opts["collection"]), nil
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
		return gopass.NewWriter(opts["store"]), nil
	case BackendInfisical:
		return infisical.NewWriter(opts["project"], opts["env"]), nil
	case BackendHashiCorpVault:
		return hashicorpvault.NewWriter(opts["addr"], opts["mount"]), nil
	case BackendSecretHub:
		return secrethub.NewWriter(opts["org"], opts["repo"]), nil
	case BackendSSMParameterStore:
		return ssmparameterstore.NewWriter(opts["region"], opts["prefix"]), nil
	case BackendAkeyless:
		return akeyless.NewWriter(), nil
	case BackendDryRun:
		return NewDryRunWriter(), nil
	default:
		return nil, fmt.Errorf("no writer registered for backend %q", b)
	}
}
