package infisical

import (
	"fmt"
	"strings"
	"testing"
)

func fakeExecSuccess(_ string, _ ...string) ([]byte, error) {
	return []byte("Secret set successfully."), nil
}

func fakeExecFailure(_ string, _ ...string) ([]byte, error) {
	return []byte("error: unauthorized"), fmt.Errorf("exit status 1")
}

func captureArgs(captured *[]string) func(string, ...string) ([]byte, error) {
	return func(name string, args ...string) ([]byte, error) {
		*captured = append([]string{name}, args...)
		return []byte("ok"), nil
	}
}

func TestWriteSecret_Success(t *testing.T) {
	w := NewWriter("proj-123", "dev")
	w.execCmd = fakeExecSuccess

	if err := w.WriteSecret("myapp", "API_KEY", "supersecret"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestWriteSecret_Failure(t *testing.T) {
	w := NewWriter("proj-123", "dev")
	w.execCmd = fakeExecFailure

	err := w.WriteSecret("myapp", "API_KEY", "supersecret")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "infisical") {
		t.Errorf("expected error to mention 'infisical', got: %v", err)
	}
}

func TestWriteSecret_PassesCorrectArgs(t *testing.T) {
	w := NewWriter("proj-abc", "prod")
	var captured []string
	w.execCmd = captureArgs(&captured)

	if err := w.WriteSecret("payments", "DB_PASS", "s3cr3t"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectContains := map[string]bool{
		"infisical":       false,
		"secrets":         false,
		"set":             false,
		"DB_PASS=s3cr3t":  false,
		"--projectId":     false,
		"proj-abc":        false,
		"--env":           false,
		"prod":            false,
		"--path":          false,
		"/payments":       false,
	}
	for _, arg := range captured {
		if _, ok := expectContains[arg]; ok {
			expectContains[arg] = true
		}
	}
	for arg, found := range expectContains {
		if !found {
			t.Errorf("expected argument %q not found in: %v", arg, captured)
		}
	}
}

func TestWriteSecret_NamespaceAsPath(t *testing.T) {
	w := NewWriter("proj-xyz", "staging")
	var captured []string
	w.execCmd = captureArgs(&captured)

	if err := w.WriteSecret("backend", "TOKEN", "abc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	pathFound := false
	for i, arg := range captured {
		if arg == "--path" && i+1 < len(captured) && captured[i+1] == "/backend" {
			pathFound = true
			break
		}
	}
	if !pathFound {
		t.Errorf("expected --path /backend in args: %v", captured)
	}
}
