package migrate

import (
	"fmt"
	"io"
)

// DryRunWriter is a SecretWriter that prints secrets instead of writing them.
type DryRunWriter struct {
	out io.Writer
}

// NewDryRunWriter returns a DryRunWriter that writes output to w.
func NewDryRunWriter(w io.Writer) *DryRunWriter {
	return &DryRunWriter{out: w}
}

// WriteSecret prints the namespace, key, and a masked value to the writer.
func (d *DryRunWriter) WriteSecret(namespace, key, value string) error {
	masked := maskValue(value)
	_, err := fmt.Fprintf(d.out, "[dry-run] namespace=%q key=%q value=%s\n", namespace, key, masked)
	return err
}

// maskValue returns a partially masked version of the secret value.
// It reveals at most the first 2 characters and replaces the rest with '*'.
func maskValue(value string) string {
	if len(value) == 0 {
		return "(empty)"
	}
	const visible = 2
	if len(value) <= visible {
		return "**"
	}
	mask := make([]byte, len(value)-visible)
	for i := range mask {
		mask[i] = '*'
	}
	return value[:visible] + string(mask)
}
