package vault

import (
	"errors"
	"strings"
	"testing"
)

func fakeExecSuccess(_ string, _ ...string) ([]byte, error) {
	return []byte("Success! Data written to: secret/data/envchain/myapp"), nil
}

func fakeExecFailure(_ string, _ ...string) ([]byte, error) {
	return []byte("Error writing data: permission denied"), errors.New("exit status 2")
}

var capturedArgs []string

func captureArgs(_ string, args ...string) ([]byte, error) {
	capturedArgs = args
	return []byte("Success!"), nil
}

func TestWriteSecret_Success(t *testing.T) {
	w := NewWriter("secret", "envchain")
	w.execCmd = fakeExecSuccess

	if err := w.WriteSecret("myapp", "API_KEY", "abc123"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestWriteSecret_Failure(t *testing.T) {
	w := NewWriter("secret", "envchain")
	w.execCmd = fakeExecFailure

	err := w.WriteSecret("myapp", "API_KEY", "abc123")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "vault kv put failed") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestWriteSecret_PassesCorrectArgs(t *testing.T) {
	w := NewWriter("secret", "envchain")
	w.execCmd = captureArgs

	if err := w.WriteSecret("myapp", "DB_PASS", "s3cr3t"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedPath := "secret/envchain/myapp"
	if len(capturedArgs) < 3 || capturedArgs[2] != expectedPath {
		t.Errorf("expected path %q, got args: %v", expectedPath, capturedArgs)
	}

	expectedKV := "DB_PASS=s3cr3t"
	if len(capturedArgs) < 4 || capturedArgs[3] != expectedKV {
		t.Errorf("expected kv %q, got args: %v", expectedKV, capturedArgs)
	}
}

func TestBackend(t *testing.T) {
	w := NewWriter("secret", "envchain")
	if w.Backend() != "vault" {
		t.Errorf("expected backend 'vault', got %q", w.Backend())
	}
}
