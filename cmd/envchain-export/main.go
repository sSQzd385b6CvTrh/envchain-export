package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/yourorg/envchain-export/internal/bitwarden"
	"github.com/yourorg/envchain-export/internal/doppler"
	"github.com/yourorg/envchain-export/internal/envchain"
	"github.com/yourorg/envchain-export/internal/migrate"
	"github.com/yourorg/envchain-export/internal/onepassword"
	"github.com/yourorg/envchain-export/internal/vault"
)

func main() {
	backendFlag := flag.String("backend", "dry-run", fmt.Sprintf("target backend (%s)", joinBackends()))
	vaultMount := flag.String("vault-mount", "secret", "Vault KV mount path")
	opVault := flag.String("op-vault", "Private", "1Password vault name")
	bwOrg := flag.String("bw-org", "", "Bitwarden organization ID (optional)")
	bwCollection := flag.String("bw-collection", "", "Bitwarden collection ID (optional)")
	dopplerProject := flag.String("doppler-project", "", "Doppler project name")
	dopplerConfig := flag.String("doppler-config", "dev", "Doppler config name")
	flag.Parse()

	backend, err := migrate.ParseBackend(*backendFlag)
	if err != nil {
		log.Fatalf("invalid backend: %v", err)
	}

	reader := envchain.NewReader()

	var writer migrate.SecretWriter
	switch backend {
	case migrate.Backend1Password:
		writer = onepassword.NewWriter(*opVault)
	case migrate.BackendDoppler:
		writer = doppler.NewWriter(*dopplerProject, *dopplerConfig)
	case migrate.BackendVault:
		writer = vault.NewWriter(*vaultMount)
	case migrate.BackendBitwarden:
		writer = bitwarden.NewWriter(*bwOrg, *bwCollection)
	case migrate.BackendDryRun:
		writer = migrate.NewDryRunWriter(os.Stdout)
	}

	m := migrate.New(reader, writer)
	summary, err := m.MigrateAll()
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	fmt.Printf("Migration complete: %d succeeded, %d failed\n", summary.Succeeded, summary.Failed)
}

func joinBackends() string {
	return fmt.Sprintf("%v", migrate.KnownBackends)
}
