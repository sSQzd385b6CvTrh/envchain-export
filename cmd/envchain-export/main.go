package main

import (
	"flag"
	"log"
	"os"

	"github.com/yourorg/envchain-export/internal/envchain"
	"github.com/yourorg/envchain-export/internal/migrate"
	"github.com/yourorg/envchain-export/internal/onepassword"
)

func main() {
	vault := flag.String("vault", "", "1Password vault name (required)")
	envchainBin := flag.String("envchain", "envchain", "path to envchain binary")
	flag.Parse()

	if *vault == "" {
		log.Println("error: --vault is required")
		flag.Usage()
		os.Exit(1)
	}

	reader := envchain.NewReader(*envchainBin)
	writer := onepassword.NewWriter(*vault)
	migrator := migrate.New(reader, writer)

	log.Printf("starting migration to 1Password vault %q", *vault)

	results, err := migrator.MigrateAll()
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	migrate.Summary(results)

	for _, r := range results {
		if r.Err != nil {
			os.Exit(1)
		}
	}
}
