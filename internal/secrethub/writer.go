package secrethub

import (
	"fmt"
	"os/exec"
	"strings"
)

// Writer writes secrets to SecretHub (via the secrethub CLI).
type Writer struct {
	execCmd func(name string, args ...string) ([]byte, error)
	repo    string
}

func defaultExecCmd(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).CombinedOutput()
}

// NewWriter creates a new SecretHub Writer.
// repo should be in the form "owner/repo".
func NewWriter(repo string) *Writer {
	return &Writer{
		execCmd: defaultExecCmd,
		repo:    repo,
	}
}

// WriteSecret writes a single key/value secret to SecretHub under the given namespace.
// The secret path will be: <repo>/<namespace>/<key>
func (w *Writer) WriteSecret(namespace, key, value string) error {
	path := fmt.Sprintf("%s/%s/%s", w.repo, namespace, key)

	// First ensure the directory exists
	dirPath := fmt.Sprintf("%s/%s", w.repo, namespace)
	out, err := w.execCmd("secrethub", "mkdir", "--parents", dirPath)
	if err != nil {
		outStr := strings.TrimSpace(string(out))
		// Ignore error if directory already exists
		if !strings.Contains(outStr, "already exists") {
			return fmt.Errorf("secrethub mkdir %s: %w: %s", dirPath, err, outStr)
		}
	}

	out, err = w.execCmd("secrethub", "write", "--in-file", "-", path)
	_ = out // write reads from stdin; we pass value via a different approach

	// Use echo-pipe approach via sh -c for value injection
	cmd := fmt.Sprintf("printf '%%s' %s | secrethub write %s", shellQuote(value), shellQuote(path))
	out, err = w.execCmd("sh", "-c", cmd)
	if err != nil {
		return fmt.Errorf("secrethub write %s: %w: %s", path, err, strings.TrimSpace(string(out)))
	}
	return nil
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\'')
	 + "'"
}
