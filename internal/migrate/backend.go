package migrate

import (
	"fmt"
	"strings"

	"github.com/envchain-export/internal/age"
	"github.com/envchain-export/internal/akeyless"
	"github.com/envchain-export/internal/awssecretsmanager"
	"github.com/envchain-export/internal/azurekeyvault"
	"github.com/envchain-export/internal/bitwarden"
	"github.com/envchain-export/internal/doppler"
	"github.com/envchain-export/internal/gcpsecretmanager"
	"github.com/envchain-export/internal/gopass"
	"github.com/envchain-export/internal/hashicorpvault"
	"github.com/envchain-export/internal/infisical"
	"github.com/envchain-export/internal/keepass"
	"github.com/envchain-export/internal/keychain"
	"github.com/envchain-export/internal/lastpass"
	"github.com/envchain-export/internal/onepassword"
	"github.com/envchain-export/internal/passbolt"
	"github.com/envchain-export/internal/secrethub"
	"github.com/envchain-export/internal/ssmparameterstore"
	"github.com/envchain-export/internal/vault"
)

// BackendType represents a supported secret backend.
type BackendType string

const (
	Backend1Password     BackendType = "1password"
	BackendDoppler       BackendType = "doppler"
	BackendVault         BackendType = "vault"
	BackendBitwarden     BackendType = "bitwarden"
	BackendAWSSecrets    BackendType = "awssecretsmanager"
	BackendGCPSecret     BackendType = "gcpsecretmanager"
	BackendAzureKeyVault BackendType = "azurekeyvault"
	BackendLastPass      BackendType = "lastpass"
	BackendKeychain      BackendType = "keychain"
	BackendGopass        BackendType = "gopass"
	BackendInfisical     BackendType = "infisical"
	BackendHashiCorpVault BackendType = "hashicorpvault"
	BackendSecretHub     BackendType = "secrethub"
	BackendSSMParameter  BackendType = "ssmparameterstore"
	BackendAkeyless      BackendType = "akeyless"
	BackendKeePass       BackendType = "keepass"
	BackendAge           BackendType = "age"
	BackendPassbolt      BackendType = "passbolt"
	BackendDryRun        BackendType = "dryrun"
)

var knownBackends = []BackendType{
	Backend1Password, BackendDoppler, BackendVault, BackendBitwarden,
	BackendAWSSecrets, BackendGCPSecret, BackendAzureKeyVault, BackendLastPass,
	BackendKeychain, BackendGopass, BackendInfisical, BackendHashiCorpVault,
	BackendSecretHub, BackendSSMParameter, BackendAkeyless, BackendKeePass,
	BackendAge, BackendPassbolt, BackendDryRun,
}

func (b BackendType) String() string { return string(b) }

// joinBackendNames returns a comma-separated list of known backend names.
func joinBackendNames() string {
	names := make([]string, len(knownBackends))
	for i, b := range knownBackends {
		names[i] = b.String()
	}
	return strings.Join(names, ", ")
}

// ParseBackend parses a string into a BackendType.
func ParseBackend(s string) (BackendType, error) {
	for _, b := range knownBackends {
		if strings.EqualFold(s, b.String()) {
			return b, nil
		}
	}
	return "", fmt.Errorf("unknown backend %q; known backends: %s", s, joinBackendNames())
}

// SecretWriter is the interface implemented by all backend writers.
type SecretWriter interface {
	WriteSecret(namespace, key, value string) error
}

// NewBackendWriter constructs the appropriate SecretWriter for the given backend.
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
	case BackendAWSSecrets:
		return awssecretsmanager.NewWriter(opts["region"]), nil
	case BackendGCPSecret:
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
		return hashicorpvault.NewWriter(opts["addr"], opts["path"]), nil
	case BackendSecretHub:
		return secrethub.NewWriter(opts["org"], opts["repo"]), nil
	case BackendSSMParameter:
		return ssmparameterstore.NewWriter(opts["region"], opts["prefix"]), nil
	case BackendAkeyless:
		return akeyless.NewWriter(opts["path"]), nil
	case BackendKeePass:
		return keepass.NewWriter(opts["db"], opts["group"]), nil
	case BackendAge:
		return age.NewWriter(opts["recipient"], opts["output"]), nil
	case BackendPassbolt:
		return passbolt.NewWriter(opts["group"]), nil
	case BackendDryRun:
		return NewDryRunWriter(), nil
	default:
		return nil, fmt.Errorf("no writer registered for backend %q", b)
	}
}
