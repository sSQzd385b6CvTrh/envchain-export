package vault

import (
	"fmt"
	"os/exec"
)

// Writer writes secrets to HashiCorp Vault via the vault CLI.
type Writer struct {
	execCmd func(name string, args ...string) ([]byte, error)
	mount   string
	path    string
}

// NewWriter creates a new Vault writer targeting the given mount and path prefix.
func NewWriter(mount, path string) *Writer {
	return &Writer{
		execCmd: defaultExecCmd,
		mount:   mount,
		path:    path,
	}
}

func defaultExecCmd(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).CombinedOutput()
}

// WriteSecret writes a single key/value pair into Vault under the configured
// mount and path, using the namespace as a sub-path.
func (w *Writer) WriteSecret(namespace, key, value string) error {
	secretPath := fmt.Sprintf("%s/%s/%s", w.mount, w.path, namespace)
	kv := fmt.Sprintf("%s=%s", key, value)

	out, err := w.execCmd("vault", "kv", "put", secretPath, kv)
	if err != nil {
		return fmt.Errorf("vault kv put failed for %s/%s: %w\noutput: %s", namespace, key, err, out)
	}
	return nil
}

// Backend returns the string identifier for this writer.
func (w *Writer) Backend() string {
	return "vault"
}
