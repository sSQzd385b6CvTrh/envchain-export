package bitwarden

import (
	"fmt"
	"os/exec"
)

// Writer writes secrets to Bitwarden via the bw CLI.
type Writer struct {
	execCmd func(name string, args ...string) ([]byte, error)
	organization string
	collection   string
}

func defaultExecCmd(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).CombinedOutput()
}

// NewWriter creates a new Bitwarden Writer.
// organization and collection are optional (empty string to skip).
func NewWriter(organization, collection string) *Writer {
	return &Writer{
		execCmd:      defaultExecCmd,
		organization: organization,
		collection:   collection,
	}
}

// WriteSecret writes a single key/value secret under the given namespace
// by creating a Bitwarden item using the bw CLI.
func (w *Writer) WriteSecret(namespace, key, value string) error {
	itemName := fmt.Sprintf("%s/%s", namespace, key)

	// Build JSON payload for bw create item
	payload := fmt.Sprintf(
		`{"type":2,"name":%q,"login":null,"secureNote":{"type":0},"notes":%q}`,
		itemName, value,
	)

	args := []string{"create", "item", payload}

	if w.organization != "" {
		args = append(args, "--organizationid", w.organization)
	}

	if w.collection != "" {
		args = append(args, "--collectionids", w.collection)
	}

	out, err := w.execCmd("bw", args...)
	if err != nil {
		return fmt.Errorf("bitwarden: failed to write secret %q: %w (output: %s)", itemName, err, string(out))
	}

	return nil
}
