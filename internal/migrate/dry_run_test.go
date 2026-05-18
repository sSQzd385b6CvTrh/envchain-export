package migrate

import (
	"bytes"
	"strings"
	"testing"
)

func TestDryRunWriter_WriteSecret_Output(t *testing.T) {
	var buf bytes.Buffer
	w := NewDryRunWriter(&buf)

	err := w.WriteSecret("myapp", "API_KEY", "supersecret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "[dry-run]") {
		t.Errorf("expected [dry-run] prefix, got: %q", got)
	}
	if !strings.Contains(got, `namespace="myapp"`) {
		t.Errorf("expected namespace in output, got: %q", got)
	}
	if !strings.Contains(got, `key="API_KEY"`) {
		t.Errorf("expected key in output, got: %q", got)
	}
	if strings.Contains(got, "supersecret") {
		t.Errorf("secret value should be masked, got: %q", got)
	}
}

func TestMaskValue_EmptyString(t *testing.T) {
	if got := maskValue(""); got != "(empty)" {
		t.Errorf("expected (empty), got %q", got)
	}
}

func TestMaskValue_ShortString(t *testing.T) {
	if got := maskValue("x"); got != "**" {
		t.Errorf("expected **, got %q", got)
	}
}

func TestMaskValue_NormalString(t *testing.T) {
	got := maskValue("abcdef")
	if !strings.HasPrefix(got, "ab") {
		t.Errorf("expected prefix 'ab', got %q", got)
	}
	if strings.Contains(got, "cdef") {
		t.Errorf("tail should be masked, got %q", got)
	}
	if len(got) != len("abcdef") {
		t.Errorf("masked value should preserve length, got %q", got)
	}
}

func TestDryRunWriter_WriteSecret_MultipleSecrets(t *testing.T) {
	var buf bytes.Buffer
	w := NewDryRunWriter(&buf)

	secrets := []struct{ ns, key, val string }{
		{"app1", "DB_PASS", "hunter2"},
		{"app2", "TOKEN", "abc123xyz"},
	}
	for _, s := range secrets {
		if err := w.WriteSecret(s.ns, s.key, s.val); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d: %q", len(lines), buf.String())
	}
}
