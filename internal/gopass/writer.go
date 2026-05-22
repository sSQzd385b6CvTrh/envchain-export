package gopass

import (
	"fmt"
	"os/exec"
	"strings"
)

// Writer writes secrets to gopass using the gopass CLI.
type Writer struct {
	execCmd func(name string, args ...string) ([]byte, error)
	store  string
}

func defaultExecCmd(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).CombinedOutput()
}

// NewWriter creates a new gopass Writer.
// store is an optional gopass mount/store prefix (e.g. "work"). Pass "" to use the default store.
func NewWriter(store string) *Writer {
	return &Writer{
		execCmd: defaultExecCmd,
		store:   store,
	}
}

// WriteSecret writes a single key/value secret under the given namespace into gopass.
// The secret path will be: [store/]namespace/key
func (w *Writer) WriteSecret(namespace, key, value string) error {
	path := w.buildPath(namespace, key)

	// gopass insert -f <path> accepts the secret value via stdin, but the CLI
	// also supports `gopass insert -f <path> <value>` in non-interactive mode.
	// We use `gopass insert --force` and pipe via echo to avoid a PTY requirement.
	cmd := exec.Command("gopass", "insert", "--force", path)
	cmd.Stdin = strings.NewReader(value + "\n")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("gopass insert failed for %q: %w: %s", path, err, strings.TrimSpace(string(out)))
	}
	return nil
}

func (w *Writer) buildPath(namespace, key string) string {
	if w.store != "" {
		return fmt.Sprintf("%s/%s/%s", w.store, namespace, key)
	}
	return fmt.Sprintf("%s/%s", namespace, key)
}
