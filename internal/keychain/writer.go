package keychain

import (
	"fmt"
	"os/exec"
)

// Writer writes secrets to macOS Keychain via the `security` CLI.
type Writer struct {
	execCmd func(name string, args ...string) *exec.Cmd
}

// NewWriter returns a new Keychain Writer.
func NewWriter() *Writer {
	return &Writer{execCmd: defaultExecCmd}
}

func defaultExecCmd(name string, args ...string) *exec.Cmd {
	return exec.Command(name, args...)
}

// WriteSecret stores a secret in the macOS Keychain under the given service and key.
func (w *Writer) WriteSecret(namespace, key, value string) error {
	cmd := w.execCmd(
		"security",
		"add-generic-password",
		"-s", namespace,
		"-a", key,
		"-w", value,
		"-U", // update if exists
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("keychain: failed to write secret %s/%s: %w: %s", namespace, key, err, string(out))
	}
	return nil
}
