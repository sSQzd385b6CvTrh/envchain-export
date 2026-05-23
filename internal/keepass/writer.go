package keepass

import (
	"fmt"
	"os/exec"
	"strings"
)

// Writer writes secrets to a KeePass database via the keepassxc-cli tool.
type Writer struct {
	database string
	group    string
	execCmd  func(name string, args ...string) *exec.Cmd
}

func defaultExecCmd(name string, args ...string) *exec.Cmd {
	return exec.Command(name, args...)
}

// NewWriter creates a new KeePass writer targeting the given database file and group.
func NewWriter(database, group string) *Writer {
	return &Writer{
		database: database,
		group:    group,
		execCmd:  defaultExecCmd,
	}
}

// WriteSecret writes a single key/value secret into the KeePass database.
// It uses keepassxc-cli to add or update an entry under the configured group.
func (w *Writer) WriteSecret(namespace, key, value string) error {
	entryTitle := fmt.Sprintf("%s/%s", namespace, key)
	if w.group != "" {
		entryTitle = fmt.Sprintf("%s/%s/%s", w.group, namespace, key)
	}

	// Try to add the entry; if it already exists, edit it instead.
	addCmd := w.execCmd(
		"keepassxc-cli", "add",
		"--quiet",
		"--password", value,
		w.database,
		entryTitle,
	)
	if out, err := addCmd.CombinedOutput(); err != nil {
		outStr := strings.TrimSpace(string(out))
		// keepassxc-cli exits non-zero with "already exists" when the entry is present.
		if !strings.Contains(outStr, "already exists") {
			return fmt.Errorf("keepassxc-cli add %q: %w: %s", entryTitle, err, outStr)
		}

		// Entry exists — update the password.
		editCmd := w.execCmd(
			"keepassxc-cli", "edit",
			"--quiet",
			"--password", value,
			w.database,
			entryTitle,
		)
		if out2, err2 := editCmd.CombinedOutput(); err2 != nil {
			return fmt.Errorf("keepassxc-cli edit %q: %w: %s", entryTitle, err2, strings.TrimSpace(string(out2)))
		}
	}

	return nil
}
