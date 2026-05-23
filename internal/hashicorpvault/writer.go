package hashicorpvault

import (
	"fmt"
	"os/exec"
)

// Writer writes secrets to HashiCorp Vault using the vault CLI.
type Writer struct {
	mount   string
	execCmd func(name string, args ...string) ([]byte, error)
}

func defaultExecCmd(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).CombinedOutput()
}

// NewWriter creates a new HashiCorp Vault writer.
// mount is the KV secrets engine mount path (e.g. "secret").
func NewWriter(mount string) *Writer {
	if mount == "" {
		mount = "secret"
	}
	return &Writer{
		mount:   mount,
		execCmd: defaultExecCmd,
	}
}

// WriteSecret writes a key/value secret to HashiCorp Vault under the given namespace.
// It uses `vault kv put <mount>/<namespace> <key>=<value>`.
func (w *Writer) WriteSecret(namespace, key, value string) error {
	path := fmt.Sprintf("%s/%s", w.mount, namespace)
	kvPair := fmt.Sprintf("%s=%s", key, value)

	out, err := w.execCmd("vault", "kv", "put", path, kvPair)
	if err != nil {
		return fmt.Errorf("hashicorpvault: failed to write secret %q to %q: %w\noutput: %s",
			key, path, err, string(out))
	}
	return nil
}
