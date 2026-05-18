package migrate

import (
	"fmt"
	"log"
)

// SecretReader can list namespaces and read secrets from a source backend.
type SecretReader interface {
	ListNamespaces() ([]string, error)
	ReadSecrets(namespace string) (map[string]string, error)
}

// SecretWriter can persist a namespace and its secrets to a target backend.
type SecretWriter interface {
	EnsureItem(namespace string, secrets map[string]string) error
}

// Result holds the outcome of migrating a single namespace.
type Result struct {
	Namespace string
	Err       error
}

// Migrator orchestrates reading from a source and writing to a destination.
type Migrator struct {
	reader SecretReader
	writer SecretWriter
}

// New creates a new Migrator.
func New(reader SecretReader, writer SecretWriter) *Migrator {
	return &Migrator{reader: reader, writer: writer}
}

// MigrateAll reads every namespace from the source and writes it to the
// destination. It returns a slice of Results, one per namespace.
func (m *Migrator) MigrateAll() ([]Result, error) {
	namespaces, err := m.reader.ListNamespaces()
	if err != nil {
		return nil, fmt.Errorf("listing namespaces: %w", err)
	}

	results := make([]Result, 0, len(namespaces))
	for _, ns := range namespaces {
		log.Printf("migrating namespace: %s", ns)
		r := Result{Namespace: ns}

		secrets, err := m.reader.ReadSecrets(ns)
		if err != nil {
			r.Err = fmt.Errorf("reading secrets: %w", err)
			results = append(results, r)
			continue
		}

		if err := m.writer.EnsureItem(ns, secrets); err != nil {
			r.Err = fmt.Errorf("writing secrets: %w", err)
		}
		results = append(results, r)
	}
	return results, nil
}

// Summary prints a human-readable migration summary to stdout.
func Summary(results []Result) {
	ok, failed := 0, 0
	for _, r := range results {
		if r.Err != nil {
			log.Printf("FAIL [%s]: %v", r.Namespace, r.Err)
			failed++
		} else {
			log.Printf("OK   [%s]", r.Namespace)
			ok++
		}
	}
	log.Printf("migration complete: %d succeeded, %d failed", ok, failed)
}
