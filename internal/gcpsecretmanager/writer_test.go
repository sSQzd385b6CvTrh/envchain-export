package gcpsecretmanager

import (
	"fmt"
	"strings"
	"testing"
)

func fakeExecSuccess(_ string, args ...string) ([]byte, error) {
	return []byte("ok"), nil
}

func fakeExecFailure(_ string, args ...string) ([]byte, error) {
	return []byte("error: something went wrong"), fmt.Errorf("exit status 1")
}

func captureArgs(captured *[]string) func(string, ...string) ([]byte, error) {
	return func(name string, args ...string) ([]byte, error) {
		*captured = append([]string{name}, args...)
		return []byte("ok"), nil
	}
}

func TestWriteSecret_Success(t *testing.T) {
	w := NewWriter("my-project")
	w.execCmd = fakeExecSuccess

	if err := w.WriteSecret("myapp", "DB_PASSWORD", "secret123"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestWriteSecret_PassesCorrectArgs(t *testing.T) {
	var captured []string
	w := NewWriter("test-project")
	w.execCmd = captureArgs(&captured)

	_ = w.WriteSecret("ns", "KEY", "val")

	// The last call should be to sh -c with the correct secret ID.
	if len(captured) == 0 {
		t.Fatal("no exec call captured")
	}
	// Verify project and secret ID appear in the shell command.
	cmd := strings.Join(captured, " ")
	if !strings.Contains(cmd, "ns_KEY") {
		t.Errorf("expected secret ID ns_KEY in command, got: %s", cmd)
	}
	if !strings.Contains(cmd, "test-project") {
		t.Errorf("expected project test-project in command, got: %s", cmd)
	}
}

func TestWriteSecret_Failure(t *testing.T) {
	w := NewWriter("my-project")
	w.execCmd = fakeExecFailure

	err := w.WriteSecret("ns", "KEY", "value")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "gcpsecretmanager") {
		t.Errorf("expected error to mention package, got: %v", err)
	}
}

func TestShellQuote(t *testing.T) {
	cases := []struct {
		input string
		want string
	}{
		{"simple", "'simple'"},
		{"with'quote", "'with'\\''quote'"},
		{"", "''"},
	}
	for _, tc := range cases {
		got := shellQuote(tc.input)
		if got != tc.want {
			t.Errorf("shellQuote(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}
