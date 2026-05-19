package gcpsecretmanager

import (
	"fmt"
	"os/exec"
	"strings"
)

// Writer writes secrets to GCP Secret Manager using the gcloud CLI.
type Writer struct {
	project string
	execCmd func(name string, args ...string) ([]byte, error)
}

func defaultExecCmd(name string, args ...string) ([]byte, error) {
	out, err := exec.Command(name, args...).CombinedOutput()
	return out, err
}

// NewWriter creates a new GCP Secret Manager writer for the given project.
func NewWriter(project string) *Writer {
	return &Writer{
		project: project,
		execCmd: defaultExecCmd,
	}
}

// WriteSecret writes a single key/value secret to GCP Secret Manager.
// It creates the secret if it does not exist, then adds a new version.
func (w *Writer) WriteSecret(namespace, key, value string) error {
	secretID := fmt.Sprintf("%s_%s", namespace, key)

	// Attempt to create the secret (idempotent; ignore already-exists error).
	_, err := w.execCmd("gcloud", "secrets", "create", secretID,
		"--project", w.project,
		"--replication-policy", "automatic",
	)
	if err != nil {
		// gcloud returns non-zero if the secret already exists; we tolerate that.
		// Any other failure will surface when we try to add the version below.
	}

	// Add a new secret version with the provided value.
	out, err := w.execCmd("gcloud", "secrets", "versions", "add", secretID,
		"--project", w.project,
		"--data-file", "-",
		"--payload-file-path", "-",
	)
	// gcloud reads the secret value from stdin; we pass it via a wrapper below.
	_ = out
	if err != nil {
		return w.addSecretVersion(secretID, value)
	}
	return nil
}

// addSecretVersion shells out to gcloud, piping value via echo.
func (w *Writer) addSecretVersion(secretID, value string) error {
	// Use printf | gcloud via sh -c to pipe the secret value.
	shCmd := fmt.Sprintf(
		"printf '%%s' %s | gcloud secrets versions add %s --project %s --data-file=-",
		shellQuote(value), secretID, w.project,
	)
	out, err := w.execCmd("sh", "-c", shCmd)
	if err != nil {
		return fmt.Errorf("gcpsecretmanager: failed to write secret %q: %w (output: %s)", secretID, err, strings.TrimSpace(string(out)))
	}
	return nil
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''" ) + "'"
}
