package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/yourorg/envchain-export/internal/doppler"
	"github.com/yourorg/envchain-export/internal/envchain"
	"github.com/yourorg/envchain-export/internal/migrate"
	"github.com/yourorg/envchain-export/internal/onepassword"
	"github.com/yourorg/envchain-export/internal/vault"
)

func main() {
	backendFlag := flag.String("backend", "", "Target backend: 1password, doppler, vault")
	vaultMount := flag.String("vault-mount", "secret", "Vault KV mount (vault backend only)")
	vaultPath := flag.String("vault-path", "envchain", "Vault KV path prefix (vault backend only)")
	opVault := flag.String("op-vault", "Private", "1Password vault name (1password backend only)")
	dopplerProject := flag.String("doppler-project", "", "Doppler project (doppler backend only)")
	dopplerConfig := flag.String("doppler-config", "dev", "Doppler config (doppler backend only)")
	flag.Parse()

	backendType, err := migrate.ParseBackend(*backendFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	reader := envchain.NewReader()

	var writer migrate.SecretWriter
	switch backendType {
	case migrate.Backend1Password:
		writer = onepassword.NewWriter(*opVault)
	case migrate.BackendDoppler:
		writer = doppler.NewWriter(*dopplerProject, *dopplerConfig)
	case migrate.BackendVault:
		writer = vault.NewWriter(*vaultMount, *vaultPath)
	}

	m := migrate.New(reader, writer)
	summary, err := m.MigrateAll()
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	fmt.Printf("Migration complete — succeeded: %d, skipped: %d, failed: %d\n",
		summary.Succeeded, summary.Skipped, summary.Failed)
}
