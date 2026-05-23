package secrethub

import (
	"fmt"
	"strings"
	"testing"
)

func fakeExecSuccess(name string, args ...string) ([]byte, error) {
	return []byte("ok"), nil
}

func fakeExecFailure(name string, args ...string) ([]byte, error) {
	return []byte("error: something went wrong"), fmt.Errorf("exit status 1")
}

func captureArgs(captured *[]string) func(string, ...string) ([]byte, error) {
	return func(name string, args ...string) ([]byte, error) {
		*captured = append(*captured, name)
		*captured = append(*captured, args...)
		return []byte("ok"), nil
	}
}

func TestWriteSecret_Success(t *testing.T) {
	w := NewWriter("myorg/myrepo")
	w.execCmd = fakeExecSuccess

	if err := w.WriteSecret("myapp", "DB_PASSWORD", "s3cr3t"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestWriteSecret_Failure(t *testing.T) {
	w := NewWriter("myorg/myrepo")
	w.execCmd = fakeExecFailure

	err := w.WriteSecret("myapp", "DB_PASSWORD", "s3cr3t")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "secrethub") {
		t.Errorf("expected error to mention 'secrethub', got: %v", err)
	}
}

func TestWriteSecret_PassesCorrectArgs(t *testing.T) {
	var captured []string
	w := NewWriter("myorg/myrepo")
	w.execCmd = captureArgs(&captured)

	_ = w.WriteSecret("myapp", "API_KEY", "abc123")

	found := false
	for _, arg := range captured {
		if strings.Contains(arg, "myorg/myrepo/myapp/API_KEY") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected secret path in args, got: %v", captured)
	}
}

func TestWriteSecret_DirectoryAlreadyExists(t *testing.T) {
	callCount := 0
	w := NewWriter("myorg/myrepo")
	w.execCmd = func(name string, args ...string) ([]byte, error) {
		callCount++
		// First call is mkdir — simulate already exists
		if callCount == 1 {
			return []byte("directory already exists"), fmt.Errorf("exit status 1")
		}
		return []byte("ok"), nil
	}

	if err := w.WriteSecret("myapp", "TOKEN", "value"); err != nil {
		t.Fatalf("expected no error when dir already exists, got: %v", err)
	}
}

func TestNewWriter_SetsRepo(t *testing.T) {
	w := NewWriter("acme/secrets")
	if w.repo != "acme/secrets" {
		t.Errorf("expected repo 'acme/secrets', got '%s'", w.repo)
	}
}
