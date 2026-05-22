package lastpass

import (
	"fmt"
	"os/exec"
	"strings"
)

// Writer writes secrets to LastPass using the lpass CLI.
type Writer struct {
	execCmd func(name string, args ...string) ([]byte, error)
	group  string
}

func defaultExecCmd(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).CombinedOutput()
}

// NewWriter creates a new LastPass writer.
// group is an optional folder/group prefix (e.g. "envchain").
func NewWriter(group string) *Writer {
	return &Writer{
		execCmd: defaultExecCmd,
		group:  group,
	}
}

// WriteSecret stores a secret in LastPass under "<group>/<namespace>" with
// the given key set as the note field and value as the password.
func (w *Writer) WriteSecret(namespace, key, value string) error {
	itemName := w.itemName(namespace)

	// lpass add --non-interactive --name <name> expects input via stdin;
	// we use `lpass edit` with --password to update or create entries.
	// Build a note-style entry: "<key>=<value>" stored as the password.
	entry := fmt.Sprintf("%s=%s", key, value)

	args := []string{
		"add",
		"--non-interactive",
		"--sync=now",
		fmt.Sprintf("--name=%s/%s", itemName, key),
		fmt.Sprintf("--password=%s", entry),
	}

	out, err := w.execCmd("lpass", args...)
	if err != nil {
		outStr := strings.TrimSpace(string(out))
		return fmt.Errorf("lpass add failed for %s/%s: %w: %s", namespace, key, err, outStr)
	}
	return nil
}

func (w *Writer) itemName(namespace string) string {
	if w.group != "" {
		return fmt.Sprintf("%s/%s", w.group, namespace)
	}
	return namespace
}
