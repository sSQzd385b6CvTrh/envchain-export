package conjur

import (
	"fmt"
	"os/exec"
	"strings"
)

// Writer writes secrets to CyberArk Conjur using the conjur CLI.
type Writer struct {
	account string
	execCmd func(name string, args ...string) *exec.Cmd
}

func defaultExecCmd(name string, args ...string) *exec.Cmd {
	return exec.Command(name, args...)
}

// NewWriter creates a new Conjur Writer.
// account is the Conjur account name (e.g. "myorg").
func NewWriter(account string) *Writer {
	return &Writer{
		account: account,
		execCmd: defaultExecCmd,
	}
}

// WriteSecret writes a single secret to Conjur under the path
// "<account>/variable/<namespace>/<key>".
func (w *Writer) WriteSecret(namespace, key, value string) error {
	variableID := fmt.Sprintf("%s:variable:%s/%s", w.account, namespace, key)

	cmd := w.execCmd("conjur", "variable", "set", "-i", variableID, "-v", value)
	out, err := cmd.CombinedOutput()
	if err != nil {
		msg := strings.TrimSpace(string(out))
		return fmt.Errorf("conjur: failed to set variable %q: %s", variableID, msg)
	}
	return nil
}
