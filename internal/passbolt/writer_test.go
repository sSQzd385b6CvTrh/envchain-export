package passbolt

import (
	"errors"
	"strings"
	"testing"
)

var capturedArgs []string

func fakeExecSuccess(_ string, args ...string) ([]byte, error) {
	capturedArgs = args
	return []byte("ok"), nil
}

func fakeExecFailure(_ string, args ...string) ([]byte, error) {
	capturedArgs = args
	return []byte("error: resource already exists"), errors.New("exit status 1")
}

func TestWriteSecret_Success(t *testing.T) {
	w := NewWriter("mygroup")
	w.execCmd = fakeExecSuccess

	if err := w.WriteSecret("myapp", "API_KEY", "supersecret"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestWriteSecret_PassesCorrectArgs(t *testing.T) {
	w := NewWriter("mygroup")
	w.execCmd = fakeExecSuccess
	capturedArgs = nil

	_ = w.WriteSecret("myapp", "DB_PASS", "s3cr3t")

	joined := strings.Join(capturedArgs, " ")
	if !strings.Contains(joined, "create resource") {
		t.Errorf("expected 'create resource' in args, got: %s", joined)
	}
	if !strings.Contains(joined, "myapp/DB_PASS") {
		t.Errorf("expected name 'myapp/DB_PASS' in args, got: %s", joined)
	}
	if !strings.Contains(joined, "s3cr3t") {
		t.Errorf("expected password value in args, got: %s", joined)
	}
	if !strings.Contains(joined, "mygroup") {
		t.Errorf("expected folder 'mygroup' in args, got: %s", joined)
	}
}

func TestWriteSecret_NoGroup(t *testing.T) {
	w := NewWriter("")
	w.execCmd = fakeExecSuccess
	capturedArgs = nil

	_ = w.WriteSecret("ns", "KEY", "val")

	joined := strings.Join(capturedArgs, " ")
	if strings.Contains(joined, "--folder") {
		t.Errorf("expected no --folder flag when groupName is empty, got: %s", joined)
	}
}

func TestWriteSecret_Failure(t *testing.T) {
	w := NewWriter("")
	w.execCmd = fakeExecFailure

	err := w.WriteSecret("ns", "KEY", "val")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "passbolt-cli create resource") {
		t.Errorf("unexpected error message: %v", err)
	}
}
