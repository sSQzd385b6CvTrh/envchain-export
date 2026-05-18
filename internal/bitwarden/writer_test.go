package bitwarden

import (
	"fmt"
	"strings"
	"testing"
)

var capturedArgs []string

func fakeExecSuccess(name string, args ...string) ([]byte, error) {
	capturedArgs = append([]string{name}, args...)
	return []byte(`{"id":"abc123"}`), nil
}

func fakeExecFailure(name string, args ...string) ([]byte, error) {
	capturedArgs = append([]string{name}, args...)
	return []byte("Not logged in."), fmt.Errorf("exit status 1")
}

func TestWriteSecret_Success(t *testing.T) {
	w := NewWriter("", "")
	w.execCmd = fakeExecSuccess

	err := w.WriteSecret("myapp", "API_KEY", "supersecret")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestWriteSecret_PassesCorrectArgs(t *testing.T) {
	w := NewWriter("", "")
	w.execCmd = fakeExecSuccess

	_ = w.WriteSecret("myapp", "DB_PASS", "hunter2")

	if capturedArgs[0] != "bw" {
		t.Errorf("expected command 'bw', got %q", capturedArgs[0])
	}
	if capturedArgs[1] != "create" || capturedArgs[2] != "item" {
		t.Errorf("expected args 'create item', got %v", capturedArgs[1:3])
	}
	if !strings.Contains(capturedArgs[3], "myapp/DB_PASS") {
		t.Errorf("payload missing item name, got: %s", capturedArgs[3])
	}
}

func TestWriteSecret_WithOrgAndCollection(t *testing.T) {
	w := NewWriter("org-id-123", "col-id-456")
	w.execCmd = fakeExecSuccess

	_ = w.WriteSecret("ns", "KEY", "val")

	joined := strings.Join(capturedArgs, " ")
	if !strings.Contains(joined, "--organizationid org-id-123") {
		t.Errorf("expected --organizationid flag, got: %s", joined)
	}
	if !strings.Contains(joined, "--collectionids col-id-456") {
		t.Errorf("expected --collectionids flag, got: %s", joined)
	}
}

func TestWriteSecret_Failure(t *testing.T) {
	w := NewWriter("", "")
	w.execCmd = fakeExecFailure

	err := w.WriteSecret("myapp", "SECRET", "value")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "bitwarden:") {
		t.Errorf("expected error to mention 'bitwarden:', got: %v", err)
	}
}
