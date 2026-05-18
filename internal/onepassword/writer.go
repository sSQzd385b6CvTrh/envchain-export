package onepassword

import (
	"fmt"
	"os/exec"
	"strings"
)

// Writer writes secrets to 1Password using the op CLI.
type Writer struct {
	vault   string
	execCmd func(name string, args ...string) ([]byte, error)
}

// NewWriter creates a new Writer targeting the given 1Password vault.
func NewWriter(vault string) *Writer {
	return &Writer{
		vault:   vault,
		execCmd: defaultExecCmd,
	}
}

func defaultExecCmd(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).CombinedOutput()
}

// EnsureItem creates a 1Password item for the namespace if it does not exist,
// then stores each key/value pair as a password field.
func (w *Writer) EnsureItem(namespace string, secrets map[string]string) error {
	if err := w.createItem(namespace); err != nil {
		return fmt.Errorf("creating item %q: %w", namespace, err)
	}
	for key, value := range secrets {
		if err := w.setField(namespace, key, value); err != nil {
			return fmt.Errorf("setting field %q on item %q: %w", key, namespace, err)
		}
	}
	return nil
}

func (w *Writer) createItem(namespace string) error {
	args := []string{
		"item", "create",
		"--category", "login",
		"--title", namespace,
		"--vault", w.vault,
	}
	out, err := w.execCmd("op", args...)
	if err != nil {
		// Ignore "already exists" errors.
		if strings.Contains(string(out), "already exists") {
			return nil
		}
		return fmt.Errorf("%w: %s", err, string(out))
	}
	return nil
}

func (w *Writer) setField(namespace, key, value string) error {
	field := fmt.Sprintf("%s[password]=%s", key, value)
	args := []string{
		"item", "edit", namespace,
		"--vault", w.vault,
		field,
	}
	out, err := w.execCmd("op", args...)
	if err != nil {
		return fmt.Errorf("%w: %s", err, string(out))
	}
	return nil
}
