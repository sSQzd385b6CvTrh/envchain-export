package enpass

import (
	"fmt"
	"os/exec"
)

// Writer writes secrets to Enpass via the enpass-cli tool.
type Writer struct {
	vault   string
	execCmd func(name string, args ...string) *exec.Cmd
}

func defaultExecCmd(name string, args ...string) *exec.Cmd {
	return exec.Command(name, args...)
}

// NewWriter creates a new Enpass writer.
// vault is the path or name of the Enpass vault to write into.
func NewWriter(vault string) *Writer {
	return &Writer{
		vault:   vault,
		execCmd: defaultExecCmd,
	}
}

// WriteSecret stores a secret in Enpass under the given namespace and key.
func (w *Writer) WriteSecret(namespace, key, value string) error {
	// enpass-cli add --vault <vault> --category "Login" --title <namespace/key> --password <value>
	title := fmt.Sprintf("%s/%s", namespace, key)
	cmd := w.execCmd(
		"enpass-cli",
		"add",
		"--vault", w.vault,
		"--category", "Login",
		"--title", title,
		"--password", value,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("enpass-cli add %q: %w: %s", title, err, string(out))
	}
	return nil
}
