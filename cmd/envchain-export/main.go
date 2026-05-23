package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/yourorg/envchain-export/internal/awssecretsmanager"
	"github.com/yourorg/envchain-export/internal/bitwarden"
	"github.com/yourorg/envchain-export/internal/doppler"
	"github.com/yourorg/envchain-export/internal/envchain"
	"github.com/yourorg/envchain-export/internal/gcpsecretmanager"
	"github.com/yourorg/envchain-export/internal/migrate"
	"github.com/yourorg/envchain-export/internal/onepassword"
	"github.com/yourorg/envchain-export/internal/vault"
)

func main() {
	backendFlag := flag.String("backend", "", fmt.Sprintf("Target backend (%s)", joinBackends()))
	dryRun := flag.Bool("dry-run", false, "Print secrets without writing them")
	// Backend-specific flags
	opVault := flag.String("op-vault", "", "1Password vault name")
	dopplerProject := flag.String("doppler-project", "", "Doppler project")
	dopplerConfig := flag.String("doppler-config", "", "Doppler config")
	gcpProject := flag.String("gcp-project", "", "GCP project ID for Secret Manager")
	awsRegion := flag.String("aws-region", "us-east-1", "AWS region for Secrets Manager")
	bwOrg := flag.String("bw-org", "", "Bitwarden organisation ID")
	bwCollection := flag.String("bw-collection", "", "Bitwarden collection ID")
	flag.Parse()

	if *backendFlag == "" && !*dryRun {
		log.Fatal("--backend is required (or use --dry-run)")
	}

	reader := envchain.NewReader()

	var writer migrate.SecretWriter
	if *dryRun {
		writer = migrate.NewDryRunWriter(os.Stdout)
	} else {
		b, err := migrate.ParseBackend(*backendFlag)
		if err != nil {
			log.Fatalf("invalid backend: %v\nAvailable backends: %s", err, joinBackends())
		}
		var initErr error
		writer, initErr = buildWriter(b, *opVault, *dopplerProject, *dopplerConfig, *gcpProject, *awsRegion, *bwOrg, *bwCollection)
		if initErr != nil {
			log.Fatalf("failed to initialise backend %q: %v", b, initErr)
		}
	}

	m := migrate.New(reader, writer)
	summary, err := m.MigrateAll()
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}
	fmt.Printf("Migration complete: %d succeeded, %d failed\n", summary.Succeeded, summary.Failed)
}

// buildWriter constructs the appropriate SecretWriter for the given backend,
// returning an error if required backend-specific flags are missing.
func buildWriter(b migrate.Backend, opVault, dopplerProject, dopplerConfig, gcpProject, awsRegion, bwOrg, bwCollection string) (migrate.SecretWriter, error) {
	switch b {
	case migrate.Backend1Password:
		if opVault == "" {
			return nil, fmt.Errorf("--op-vault is required for the 1password backend")
		}
		return onepassword.NewWriter(opVault), nil
	case migrate.BackendDoppler:
		if dopplerProject == "" || dopplerConfig == "" {
			return nil, fmt.Errorf("--doppler-project and --doppler-config are required for the doppler backend")
		}
		return doppler.NewWriter(dopplerProject, dopplerConfig), nil
	case migrate.BackendVault:
		return vault.NewWriter(), nil
	case migrate.BackendBitwarden:
		if bwOrg == "" || bwCollection == "" {
			return nil, fmt.Errorf("--bw-org and --bw-collection are required for the bitwarden backend")
		}
		return bitwarden.NewWriter(bwOrg, bwCollection), nil
	case migrate.BackendAWSSecretsManager:
		return awssecretsmanager.NewWriter(awsRegion), nil
	case migrate.BackendGCPSecretManager:
		if gcpProject == "" {
			return nil, fmt.Errorf("--gcp-project is required for the gcpsecretmanager backend")
		}
		return gcpsecretmanager.NewWriter(gcpProject), nil
	default:
		return nil, fmt.Errorf("unhandled backend %q", b)
	}
}

func joinBackends() string {
	names := make([]string, len(migrate.KnownBackends))
	for i, b := range migrate.KnownBackends {
		names[i] = b.String()
	}
	return strings.Join(names, ", ")
}
