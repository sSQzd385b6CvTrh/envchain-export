package envchain

import (
	"fmt"
	"os/exec"
	"strings"
)

// Secret represents a single envchain secret entry.
type Secret struct {
	Namespace string
	Key       string
	Value     string
}

// Reader reads secrets from envchain via the envchain CLI.
type Reader struct {
	envchainBin string
}

// NewReader creates a new Reader, optionally specifying the envchain binary path.
func NewReader(envchainBin string) *Reader {
	if envchainBin == "" {
		envchainBin = "envchain"
	}
	return &Reader{envchainBin: envchainBin}
}

// ListNamespaces returns all envchain namespaces stored in the keychain.
func (r *Reader) ListNamespaces() ([]string, error) {
	cmd := exec.Command(r.envchainBin, "--list")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("listing envchain namespaces: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var namespaces []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			namespaces = append(namespaces, line)
		}
	}
	return namespaces, nil
}

// ReadSecrets retrieves all environment variables for the given namespace.
func (r *Reader) ReadSecrets(namespace string) ([]Secret, error) {
	// envchain <namespace> env prints KEY=VALUE pairs for the namespace
	cmd := exec.Command(r.envchainBin, namespace, "env")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("reading secrets for namespace %q: %w", namespace, err)
	}

	var secrets []Secret
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			continue
		}
		key := line[:idx]
		value := line[idx+1:]
		secrets = append(secrets, Secret{
			Namespace: namespace,
			Key:       key,
			Value:     value,
		})
	}
	return secrets, nil
}
