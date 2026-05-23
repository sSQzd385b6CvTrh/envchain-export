package infisical

import (
	"fmt"
	"os/exec"
	"strings"
)

// Writer writes secrets to Infisical using the infisical CLI.
type Writer struct {
	projectID string
	environment string
	execCmd func(name string, args ...string) ([]byte, error)
}

func defaultExecCmd(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).CombinedOutput()
}

// NewWriter creates a new Infisical Writer.
// projectID is the Infisical project ID and environment is e.g. "dev", "prod".
func NewWriter(projectID, environment string) *Writer {
	return &Writer{
		projectID:   projectID,
		environment: environment,
		execCmd:     defaultExecCmd,
	}
}

// WriteSecret writes a single key/value secret to Infisical under the given namespace.
func (w *Writer) WriteSecret(namespace, key, value string) error {
	// Infisical secret path uses namespace as a folder path.
	path := fmt.Sprintf("/%s", namespace)

	out, err := w.execCmd(
		"infisical", "secrets", "set",
		fmt.Sprintf("%s=%s", key, value),
		"--projectId", w.projectID,
		"--env", w.environment,
		"--path", path,
		"--plain",
	)
	if err != nil {
		outStr := strings.TrimSpace(string(out))
		return fmt.Errorf("infisical: failed to set secret %s/%s: %w (output: %s)", namespace, key, err, outStr)
	}
	return nil
}
