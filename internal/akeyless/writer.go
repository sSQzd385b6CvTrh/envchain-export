package akeyless

import (
	"fmt"
	"os/exec"
)

// execCmd is the function used to run external commands (replaceable in tests).
var execCmd = defaultExecCmd

func defaultExecCmd(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).CombinedOutput()
}

// Writer writes secrets to Akeyless Vault using the `akeyless` CLI.
type Writer struct {
	execCmd func(string, ...string) ([]byte, error)
}

// NewWriter returns a new Akeyless Writer.
func NewWriter() *Writer {
	return &Writer{execCmd: execCmd}
}

// WriteSecret stores a secret in Akeyless under the path /<namespace>/<key>.
func (w *Writer) WriteSecret(namespace, key, value string) error {
	path := fmt.Sprintf("/%s/%s", namespace, key)
	out, err := w.execCmd(
		"akeyless",
		"create-secret",
		"--name", path,
		"--value", value,
	)
	if err != nil {
		return fmt.Errorf("akeyless create-secret %s: %w: %s", path, err, string(out))
	}
	return nil
}
