package ssmparameterstore

import (
	"fmt"
	"os/exec"
	"strings"
)

// Writer writes secrets to AWS SSM Parameter Store via the AWS CLI.
type Writer struct {
	execCmd func(name string, args ...string) *exec.Cmd
	path    string
	kmsKeyID string
}

func defaultExecCmd(name string, args ...string) *exec.Cmd {
	return exec.Command(name, args...)
}

// NewWriter creates a new SSM Parameter Store writer.
// path is the base path prefix (e.g. "/myapp").
// kmsKeyID is optional; pass empty string to use the default SSM key.
func NewWriter(path, kmsKeyID string) *Writer {
	return &Writer{
		execCmd:  defaultExecCmd,
		path:     strings.TrimRight(path, "/"),
		kmsKeyID: kmsKeyID,
	}
}

// WriteSecret stores a secret under <path>/<namespace>/<key> as a SecureString.
func (w *Writer) WriteSecret(namespace, key, value string) error {
	paramName := fmt.Sprintf("%s/%s/%s", w.path, namespace, key)

	args := []string{
		"ssm", "put-parameter",
		"--name", paramName,
		"--value", value,
		"--type", "SecureString",
		"--overwrite",
	}

	if w.kmsKeyID != "" {
		args = append(args, "--key-id", w.kmsKeyID)
	}

	cmd := w.execCmd("aws", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ssm put-parameter %s: %w: %s", paramName, err, strings.TrimSpace(string(out)))
	}
	return nil
}
