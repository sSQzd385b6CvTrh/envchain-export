package awssecretsmanager

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// Writer writes secrets to AWS Secrets Manager via the AWS CLI.
type Writer struct {
	execCmd func(name string, args ...string) ([]byte, error)
	region  string
}

// NewWriter creates a new AWS Secrets Manager writer.
// region may be empty to use the default AWS region from environment/config.
func NewWriter(region string) *Writer {
	return &Writer{
		execCmd: defaultExecCmd,
		region:  region,
	}
}

func defaultExecCmd(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).CombinedOutput()
}

// WriteSecret writes a single key/value pair under the given namespace (secret name).
// It attempts to create the secret first; if it already exists it updates the value.
func (w *Writer) WriteSecret(namespace, key, value string) error {
	secretID := fmt.Sprintf("%s/%s", namespace, key)

	payload, err := json.Marshal(map[string]string{key: value})
	if err != nil {
		return fmt.Errorf("awssecretsmanager: marshal payload: %w", err)
	}

	args := w.buildArgs("create-secret",
		"--name", secretID,
		"--secret-string", string(payload),
	)

	out, err := w.execCmd("aws", args...)
	if err != nil {
		if strings.Contains(string(out), "ResourceExistsException") {
			return w.updateSecret(secretID, string(payload))
		}
		return fmt.Errorf("awssecretsmanager: create-secret %s: %s", secretID, strings.TrimSpace(string(out)))
	}
	return nil
}

func (w *Writer) updateSecret(secretID, payload string) error {
	args := w.buildArgs("put-secret-value",
		"--secret-id", secretID,
		"--secret-string", payload,
	)

	out, err := w.execCmd("aws", args...)
	if err != nil {
		return fmt.Errorf("awssecretsmanager: put-secret-value %s: %s", secretID, strings.TrimSpace(string(out)))
	}
	return nil
}

func (w *Writer) buildArgs(subcommands ...string) []string {
	args := []string{"secretsmanager"}
	args = append(args, subcommands...)
	if w.region != "" {
		args = append(args, "--region", w.region)
	}
	return args
}
