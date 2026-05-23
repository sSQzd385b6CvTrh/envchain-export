package akeyless

import (
	"errors"
	"strings"
	"testing"
)

var capturedArgs []string

func fakeExecSuccess(name string, args ...string) ([]byte, error) {
	capturedArgs = append([]string{name}, args...)
	return []byte("secret created"), nil
}

func fakeExecFailure(name string, args ...string) ([]byte, error) {
	capturedArgs = append([]string{name}, args...)
	return []byte("error: unauthorized"), errors.New("exit status 1")
}

func TestWriteSecret_Success(t *testing.T) {
	w := &Writer{execCmd: fakeExecSuccess}
	if err := w.WriteSecret("myapp", "DB_PASS", "s3cr3t"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestWriteSecret_Failure(t *testing.T) {
	w := &Writer{execCmd: fakeExecFailure}
	err := w.WriteSecret("myapp", "DB_PASS", "s3cr3t")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "akeyless create-secret") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestWriteSecret_PassesCorrectArgs(t *testing.T) {
	capturedArgs = nil
	w := &Writer{execCmd: fakeExecSuccess}
	if err := w.WriteSecret("ns", "MY_KEY", "myvalue"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"akeyless", "create-secret", "--name", "/ns/MY_KEY", "--value", "myvalue"}
	if len(capturedArgs) != len(expected) {
		t.Fatalf("expected %d args, got %d: %v", len(expected), len(capturedArgs), capturedArgs)
	}
	for i, v := range expected {
		if capturedArgs[i] != v {
			t.Errorf("arg[%d]: expected %q, got %q", i, v, capturedArgs[i])
		}
	}
}

func TestWriteSecret_PathFormat(t *testing.T) {
	capturedArgs = nil
	w := &Writer{execCmd: fakeExecSuccess}
	_ = w.WriteSecret("production", "API_TOKEN", "tok")

	var path string
	for i, a := range capturedArgs {
		if a == "--name" && i+1 < len(capturedArgs) {
			path = capturedArgs[i+1]
			break
		}
	}
	if path != "/production/API_TOKEN" {
		t.Errorf("expected path /production/API_TOKEN, got %q", path)
	}
}
