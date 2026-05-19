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
			log.Fatalf("invalid backend: %v", err)
		}
		switch b {
		case migrate.Backend1Password:
			writer = onepassword.NewWriter(*opVault)
		case migrate.BackendDoppler:
			writer = doppler.NewWriter(*dopplerProject, *dopplerConfig)
		case migrate.BackendVault:
			writer = vault.NewWriter()
		case migrate.BackendBitwarden:
			writer = bitwarden.NewWriter(*bwOrg, *bwCollection)
		case migrate.BackendAWSSecretsManager:
			writer = awssecretsmanager.NewWriter(*awsRegion)
		case migrate.BackendGCPSecretManager:
			writer = gcpsecretmanager.NewWriter(*gcpProject)
		}
	}

	m := migrate.New(reader, writer)
	summary, err := m.MigrateAll()
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}
	fmt.Printf("Migration complete: %d succeeded, %d failed\n", summary.Succeeded, summary.Failed)
}

func joinBackends() string {
	names := make([]string, len(migrate.KnownBackends))
	for i, b := range migrate.KnownBackends {
		names[i] = b.String()
	}
	return strings.Join(names, ", ")
}
