package age

import (
	"fmt"
	"os/exec"
	"strings"
)

// Writer writes secrets to an age-encrypted file via the age CLI.
type Writer struct {
	recipient string
	outputFile string
	execCmd func(name string, args ...string) *exec.Cmd
}

func defaultExecCmd(name string, args ...string) *exec.Cmd {
	return exec.Command(name, args...)
}

// NewWriter creates a new age Writer.
// recipient is the age public key or recipient file path.
// outputFile is the path to the encrypted output file.
func NewWriter(recipient, outputFile string) *Writer {
	return &Writer{
		recipient:  recipient,
		outputFile: outputFile,
		execCmd:    defaultExecCmd,
	}
}

// WriteSecret encrypts a secret and appends it to the output file using age.
// The secret is stored as "namespace/key=value" in the encrypted file.
func (w *Writer) WriteSecret(namespace, key, value string) error {
	entry := fmt.Sprintf("%s/%s=%s", namespace, key, value)

	// Use age to encrypt the entry and append to the output file
	args := []string{
		"--recipient", w.recipient,
		"--output", w.outputFile,
	}

	cmd := w.execCmd("age", args...)
	cmd.Stdin = strings.NewReader(entry)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("age: failed to encrypt secret %s/%s: %w: %s", namespace, key, err, string(out))
	}

	return nil
}
