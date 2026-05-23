package hashicorpvault

import (
	"errors"
	"strings"
	"testing"
)

func fakeExecSuccess(name string, args ...string) ([]byte, error) {
	return []byte("Success! Data written to: secret/myapp"), nil
}

func fakeExecFailure(name string, args ...string) ([]byte, error) {
	return []byte("Error writing data: permission denied"), errors.New("exit status 1")
}

var capturedArgs []string

func captureArgs(name string, args ...string) ([]byte, error) {
	capturedArgs = append([]string{name}, args...)
	return []byte("Success!"), nil
}

func TestWriteSecret_Success(t *testing.T) {
	w := NewWriter("secret")
	w.execCmd = fakeExecSuccess

	if err := w.WriteSecret("myapp", "DB_PASSWORD", "s3cr3t"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestWriteSecret_Failure(t *testing.T) {
	w := NewWriter("secret")
	w.execCmd = fakeExecFailure

	err := w.WriteSecret("myapp", "DB_PASSWORD", "s3cr3t")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "hashicorpvault") {
		t.Errorf("expected error to contain 'hashicorpvault', got: %v", err)
	}
}

func TestWriteSecret_PassesCorrectArgs(t *testing.T) {
	w := NewWriter("kv")
	w.execCmd = captureArgs

	if err := w.WriteSecret("myns", "API_KEY", "abc123"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"vault", "kv", "put", "kv/myns", "API_KEY=abc123"}
	if len(capturedArgs) != len(expected) {
		t.Fatalf("expected args %v, got %v", expected, capturedArgs)
	}
	for i, arg := range expected {
		if capturedArgs[i] != arg {
			t.Errorf("arg[%d]: expected %q, got %q", i, arg, capturedArgs[i])
		}
	}
}

func TestNewWriter_DefaultMount(t *testing.T) {
	w := NewWriter("")
	if w.mount != "secret" {
		t.Errorf("expected default mount 'secret', got %q", w.mount)
	}
}

func TestNewWriter_CustomMount(t *testing.T) {
	w := NewWriter("kv")
	if w.mount != "kv" {
		t.Errorf("expected mount 'kv', got %q", w.mount)
	}
}
