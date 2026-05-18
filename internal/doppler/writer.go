package doppler

import (
	"fmt"
	"os/exec"
	"strings"
)

// Writer writes secrets to Doppler using the Doppler CLI.
type Writer struct {
	project string
	config  string
	execCmd func(name string, args ...string) ([]byte, error)
}

// NewWriter creates a new Doppler Writer for the given project and config.
func NewWriter(project, config string) *Writer {
	return &Writer{
		project: project,
		config:  config,
		execCmd: defaultExecCmd,
	}
}

func defaultExecCmd(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).CombinedOutput()
}

// WriteSecret writes a single key-value secret to Doppler.
func (w *Writer) WriteSecret(key, value string) error {
	args := []string{
		"secrets", "set",
		"--project", w.project,
		"--config", w.config,
		fmt.Sprintf("%s=%s", key, value),
	}

	out, err := w.execCmd("doppler", args...)
	if err != nil {
		return fmt.Errorf("doppler secrets set failed for key %q: %w (output: %s)", key, err, strings.TrimSpace(string(out)))
	}
	return nil
}

// WriteSecrets writes multiple key-value pairs to Doppler.
func (w *Writer) WriteSecrets(namespace string, secrets map[string]string) error {
	for key, value := range secrets {
		if err := w.WriteSecret(key, value); err != nil {
			return fmt.Errorf("namespace %q: %w", namespace, err)
		}
	}
	return nil
}
