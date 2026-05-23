package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/envchain-export/internal/envchain"
	"github.com/envchain-export/internal/migrate"
)

func main() {
	backendFlag := flag.String("backend", "", fmt.Sprintf("Target backend (%s)", joinBackends()))
	namespacesFlag := flag.String("namespaces", "", "Comma-separated list of namespaces to migrate (default: all)")
	dryRunFlag := flag.Bool("dry-run", false, "Print secrets without writing them")

	// Backend-specific options
	vaultFlag := flag.String("vault", "", "Vault/account name (1password, azurekeyvault)")
	projectFlag := flag.String("project", "", "Project name (doppler, gcpsecretmanager, infisical)")
	configFlag := flag.String("config", "", "Config/environment (doppler, infisical)")
	pathFlag := flag.String("path", "", "Path prefix (vault, hashicorpvault, akeyless, ssmparameterstore)")
	orgFlag := flag.String("org", "", "Organisation (bitwarden, secrethub)")
	collectionFlag := flag.String("collection", "", "Collection (bitwarden)")
	regionFlag := flag.String("region", "", "AWS region (awssecretsmanager, ssmparameterstore)")
	folderFlag := flag.String("folder", "", "Folder/group name (lastpass, passbolt)")
	storeFlag := flag.String("store", "", "Store name (gopass)")
	addrFlag := flag.String("addr", "", "Server address (hashicorpvault)")
	repoFlag := flag.String("repo", "", "Repository (secrethub)")
	prefixFlag := flag.String("prefix", "", "Parameter prefix (ssmparameterstore)")
	dbFlag := flag.String("db", "", "Database file path (keepass)")
	groupFlag := flag.String("group", "", "Group name (keepass)")
	recipientFlag := flag.String("recipient", "", "Age recipient public key")
	outputFlag := flag.String("output", "", "Output file (age)")

	flag.Parse()

	if *dryRunFlag {
		*backendFlag = "dryrun"
	}

	if *backendFlag == "" {
		log.Fatalf("--backend is required. Known backends: %s", joinBackends())
	}

	backend, err := migrate.ParseBackend(*backendFlag)
	if err != nil {
		log.Fatalf("invalid backend: %v", err)
	}

	opts := map[string]string{
		"vault":      *vaultFlag,
		"project":    *projectFlag,
		"config":     *configFlag,
		"path":       *pathFlag,
		"org":        *orgFlag,
		"collection": *collectionFlag,
		"region":     *regionFlag,
		"folder":     *folderFlag,
		"store":      *storeFlag,
		"addr":       *addrFlag,
		"repo":       *repoFlag,
		"prefix":     *prefixFlag,
		"db":         *dbFlag,
		"group":      *groupFlag,
		"recipient":  *recipientFlag,
		"output":     *outputFlag,
		"env":        *configFlag,
	}

	writer, err := buildWriter(backend, opts)
	if err != nil {
		log.Fatalf("failed to build writer: %v", err)
	}

	reader := envchain.NewReader()
	migrator := migrate.New(reader, writer)

	var namespaces []string
	if *namespacesFlag != "" {
		namespaces = strings.Split(*namespacesFlag, ",")
	}

	summary, err := migrator.MigrateAll(namespaces)
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	fmt.Fprintf(os.Stdout, "Migration complete: %d succeeded, %d failed\n",
		summary.Succeeded, summary.Failed)

	if summary.Failed > 0 {
		os.Exit(1)
	}
}

func buildWriter(b migrate.BackendType, opts map[string]string) (migrate.SecretWriter, error) {
	return migrate.NewBackendWriter(b, opts)
}

func joinBackends() string {
	_, err := migrate.ParseBackend("__list__")
	if err != nil {
		// Extract list from error message
		s := err.Error()
		if idx := strings.Index(s, "known backends: "); idx >= 0 {
			return s[idx+len("known backends: "):]
		}
	}
	return "see --help"
}
