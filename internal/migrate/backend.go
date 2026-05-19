package migrate

import (
	"fmt"
	"strings"

	"github.com/envchain-export/internal/awssecretsmanager"
	"github.com/envchain-export/internal/azurekeyvault"
	"github.com/envchain-export/internal/bitwarden"
	"github.com/envchain-export/internal/doppler"
	"github.com/envchain-export/internal/gcpsecretmanager"
	"github.com/envchain-export/internal/onepassword"
	"github.com/envchain-export/internal/vault"
)

// BackendType identifies a supported secret backend.
type BackendType string

const (
	Backend1Password    BackendType = "1password"
	BackendDoppler      BackendType = "doppler"
	BackendVault        BackendType = "vault"
	BackendBitwarden    BackendType = "bitwarden"
	BackendAWSSecrets   BackendType = "aws-secrets-manager"
	BackendGCPSecrets   BackendType = "gcp-secret-manager"
	BackendAzureKeyVault BackendType = "azure-keyvault"
	BackendDryRun       BackendType = "dry-run"
)

// KnownBackends lists all supported backend identifiers.
var KnownBackends = []BackendType{
	Backend1Password,
	BackendDoppler,
	BackendVault,
	BackendBitwarden,
	BackendAWSSecrets,
	BackendGCPSecrets,
	BackendAzureKeyVault,
	BackendDryRun,
}

func (b BackendType) String() string { return string(b) }

// SecretWriter is the interface all backend writers must satisfy.
type SecretWriter interface {
	WriteSecret(namespace, key, value string) error
}

// ParseBackend resolves a backend name and optional target string into a SecretWriter.
func ParseBackend(name, target string) (SecretWriter, error) {
	switch BackendType(name) {
	case Backend1Password:
		return onepassword.NewWriter(target), nil
	case BackendDoppler:
		return doppler.NewWriter(target), nil
	case BackendVault:
		return vault.NewWriter(target), nil
	case BackendBitwarden:
		return bitwarden.NewWriter(target, "", ""), nil
	case BackendAWSSecrets:
		return awssecretsmanager.NewWriter(target), nil
	case BackendGCPSecrets:
		return gcpsecretmanager.NewWriter(target), nil
	case BackendAzureKeyVault:
		return azurekeyvault.NewWriter(target), nil
	case BackendDryRun:
		return NewDryRunWriter(), nil
	default:
		return nil, fmt.Errorf("unknown backend %q: supported backends are %s", name, joinBackendNames())
	}
}

func joinBackendNames() string {
	names := make([]string, len(KnownBackends))
	for i, b := range KnownBackends {
		names[i] = b.String()
	}
	return strings.Join(names, ", ")
}
