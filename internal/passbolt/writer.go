package passbolt

import (
	"fmt"
	"os/exec"
	"strings"
)

// Writer writes secrets to Passbolt via the passbolt-cli tool.
type Writer struct {
	execCmd func(name string, args ...string) ([]byte, error)
	groupName string
}

func defaultExecCmd(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).CombinedOutput()
}

// NewWriter creates a new Passbolt Writer.
// groupName is the Passbolt group/folder to store secrets under.
func NewWriter(groupName string) *Writer {
	return &Writer{
		execCmd:   defaultExecCmd,
		groupName: groupName,
	}
}

// WriteSecret writes a single secret to Passbolt using passbolt-cli.
func (w *Writer) WriteSecret(namespace, key, value string) error {
	name := fmt.Sprintf("%s/%s", namespace, key)

	args := []string{
		"create", "resource",
		"--name", name,
		"--password", value,
	}

	if w.groupName != "" {
		args = append(args, "--folder", w.groupName)
	}

	out, err := w.execCmd("passbolt-cli", args...)
	if err != nil {
		outStr := strings.TrimSpace(string(out))
		return fmt.Errorf("passbolt-cli create resource %q: %w: %s", name, err, outStr)
	}

	return nil
}
