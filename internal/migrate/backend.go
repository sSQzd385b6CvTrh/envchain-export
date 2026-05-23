package migrate

import (
	"fmt"
	"strings"

	"github.com/nicholasgasior/envchain-export/internal/akeyless"
	"github.com/nicholasgasior/envchain-export/internal/age"
	"github.com/nicholasgasior/envchain-export/internal/awssecretsmanager"
	"github.com/nicholasgasior/envchain-export/internal/azurekeyvault"
	"github.com/nicholasgasior/envchain-export/internal/bitwarden"
	"github.com/nicholasgasior/envchain-export/internal/conjur"
	"github.com/nicholasgasior/envchain-export/internal/doppler"
	"github.com/nicholasgasior/envchain-export/internal/enpass"
	"github.com/nicholasgasior/envchain-export/internal/gcpsecretmanager"
	"github.com/nicholasgasior/envchain-export/internal/gopass"
	"github.com/nicholasgasior/envchain-export/internal/hashicorpvault"
	"github.com/nicholasgasior/envchain-export/internal/infisical"
	"github.com/nicholasgasior/envchain-export/internal/keepass"
	"github.com/nicholasgasior/envchain-export/internal/keychain"
	"github.com/nicholasgasior/envchain-export/internal/lastpass"
	"github.com/nicholasgasior/envchain-export/internal/onepassword"
	"github.com/nicholasgasior/envchain-export/internal/passbolt"
	"github.com/nicholasgasior/envchain-export/internal/secrethub"
	"github.com/nicholasgasior/envchain-export/internal/ssmparameterstore"
	"github.com/nicholasgasior/envchain-export/internal/vault"
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
	BackendInfisical
	BackendHashiCorpVault
	BackendSecretHub
	BackendSSMParameterStore
	BackendAkeyless
	BackendKeePass
	BackendAge
	BackendPassbolt
	BackendConjur
	BackendBitwarden
	BackendEnpass
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
	BackendInfisical:         "infisical",
	BackendHashiCorpVault:    "hashicorp-vault",
	BackendSecretHub:         "secrethub",
	BackendSSMParameterStore: "ssm-parameter-store",
	BackendAkeyless:          "akeyless",
	BackendKeePass:           "keepass",
	BackendAge:               "age",
	BackendPassbolt:          "passbolt",
	BackendConjur:            "conjur",
	BackendBitwarden:         "bitwarden",
	BackendEnpass:            "enpass",
}

func (b BackendType) String() string {
	if name, ok := backendNames[b]; ok {
		return name
	}
	return fmt.Sprintf("unknown(%d)", int(b))
}

func joinBackendNames() string {
	names := make([]string, 0, len(backendNames))
	for _, name := range backendNames {
		names = append(names, name)
	}
	return strings.Join(names, ", ")
}

// ParseBackend parses a backend name string into a BackendType.
func ParseBackend(name string) (BackendType, error) {
	for bt, n := range backendNames {
		if strings.EqualFold(n, name) {
			return bt, nil
		}
	}
	return 0, fmt.Errorf("unknown backend %q; supported: %s", name, joinBackendNames())
}

// SecretWriter is the interface implemented by all backend writers.
type SecretWriter interface {
	WriteSecret(namespace, key, value string) error
}

// NewBackendWriter constructs the appropriate writer for the given backend.
func NewBackendWriter(bt BackendType, flags map[string]string) (SecretWriter, error) {
	switch bt {
	case Backend1Password:
		return onepassword.NewWriter(flags["vault"]), nil
	case BackendDoppler:
		return doppler.NewWriter(flags["project"], flags["config"]), nil
	case BackendVault:
		return vault.NewWriter(flags["path"]), nil
	case BackendAWSSecretsManager:
		return awssecretsmanager.NewWriter(flags["region"]), nil
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
		return hashicorpvault.NewWriter(flags["addr"], flags["path"]), nil
	case BackendSecretHub:
		return secrethub.NewWriter(flags["org"], flags["repo"]), nil
	case BackendSSMParameterStore:
		return ssmparameterstore.NewWriter(flags["region"], flags["prefix"]), nil
	case BackendAkeyless:
		return akeyless.NewWriter(flags["path"]), nil
	case BackendKeePass:
		return keepass.NewWriter(flags["db"], flags["password"]), nil
	case BackendAge:
		return age.NewWriter(flags["recipient"], flags["output"]), nil
	case BackendPassbolt:
		return passbolt.NewWriter(flags["group"]), nil
	case BackendConjur:
		return conjur.NewWriter(flags["account"], flags["path"]), nil
	case BackendBitwarden:
		return bitwarden.NewWriter(flags["org"], flags["collection"]), nil
	case BackendEnpass:
		return enpass.NewWriter(flags["vault"]), nil
	default:
		return nil, fmt.Errorf("no writer implemented for backend %s", bt)
	}
}
