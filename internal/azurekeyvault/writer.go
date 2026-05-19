package azurekeyvault

import (
	"fmt"
	"os/exec"
	"strings"
)

// Writer writes secrets to Azure Key Vault using the az CLI.
type Writer struct {
	vaultName string
	execCmd   func(name string, args ...string) ([]byte, error)
}

func defaultExecCmd(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).CombinedOutput()
}

// NewWriter creates a new Azure Key Vault writer for the given vault name.
func NewWriter(vaultName string) *Writer {
	return &Writer{
		vaultName: vaultName,
		execCmd:   defaultExecCmd,
	}
}

// WriteSecret writes a single secret to Azure Key Vault.
// Secret names in Key Vault may only contain alphanumeric characters and dashes,
// so underscores and dots in the key are replaced with dashes.
func (w *Writer) WriteSecret(namespace, key, value string) error {
	secretName := sanitizeName(namespace + "-" + key)

	out, err := w.execCmd(
		"az", "keyvault", "secret", "set",
		"--vault-name", w.vaultName,
		"--name", secretName,
		"--value", value,
		"--output", "none",
	)
	if err != nil {
		return fmt.Errorf("az keyvault secret set %q: %w: %s", secretName, err, strings.TrimSpace(string(out)))
	}
	return nil
}

// sanitizeName replaces characters not allowed in Key Vault secret names with dashes.
func sanitizeName(s string) string {
	replacer := strings.NewReplacer(
		"_", "-",
		".", "-",
		" ", "-",
	)
	return replacer.Replace(s)
}
